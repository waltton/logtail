package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	yaml "gopkg.in/yaml.v2"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"fmt"

	"strings"

	pb "github.com/waltton/logtail/logtail"
)

var data = `
teste: /home/waltton/dev/go/src/github.com/waltton/logtail/client/main.go
teste2: /home/waltton/dev/go/src/github.com/waltton/logtail/server/main.go
`

const port = ":50051"
const maxSize = 1024

func getFiles() (map[string]string, error) {
	files := map[string]string{}

	err := yaml.Unmarshal([]byte(data), &files)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal the yaml file: %s", err)
	}

	return files, nil
}

func getFilePath(fileName string) (string, error) {
	files := map[string]string{}

	err := yaml.Unmarshal([]byte(data), &files)
	if err != nil {
		return "", fmt.Errorf("could not unmarshal the yaml file: %s", err)
	}

	path, exists := files[fileName]
	if !exists {
		return "", fmt.Errorf("file not exists")
	}

	return path, nil
}

// server is used to implement logtail.LogTailServer.
type server struct{}

// GetFiles implements logtail.LogTailServer
func (s *server) GetFiles(ctx context.Context, in *pb.RequestFile) (*pb.Files, error) {
	files, err := getFiles()
	if err != nil {
		return nil, fmt.Errorf("could get the file list: %s", err)
	}

	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}

	return &pb.Files{Name: names}, nil
}

// GetFileContent implements logtail.LogTailServer
func (s *server) GetFileContent(ctx context.Context, in *pb.FileName) (*pb.Content, error) {
	path, err := getFilePath(in.GetName())
	if err != nil {
		return nil, fmt.Errorf("could get the file path: %s", err)
	}

	fileStat, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("could read the file size: %s", err)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could open the file: %s", err)
	}

	var lines []string
	if fileStat.Size() > maxSize {
		_, err = file.Seek(maxSize, 0)
		if err != nil {
			return nil, fmt.Errorf("could seek in the file: %s", err)
		}

		buf := make([]byte, maxSize)
		_, err = file.Read(buf)
		if err != nil {
			return nil, fmt.Errorf("could read the file: %s", err)
		}

		lines = strings.Split("path: "+path+"\n\n"+string(buf), "\n")
	} else {
		dat, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("could read the file: %s", err)
		}
		lines = strings.Split("path: "+path+"\n\n"+string(dat), "\n")
	}

	return &pb.Content{Line: lines}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterLogTailServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
