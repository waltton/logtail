package main

import (
	"log"
	"net"

	yaml "gopkg.in/yaml.v2"

	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"fmt"

	pb "github.com/waltton/logtail/logtail"
)

var data = `
teste: teste
teste2: teste2
`

const port = ":50051"

func getFiles() (map[string]string, error) {
	files := map[string]string{}

	err := yaml.Unmarshal([]byte(data), &files)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal the yaml file: %s", err)
	}

	return files, nil
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
