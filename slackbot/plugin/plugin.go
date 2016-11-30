package plugin

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-chat-bot/bot"
	"google.golang.org/grpc"

	pb "github.com/waltton/logtail/logtail"
)

var servers = map[string]string{
	"teste.teste": "localhost:50051",
}

func newConnection(address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn, nil
}

func serverList() []string {
	result := make([]string, 0, len(servers))
	for key := range servers {
		result = append(result, key)
	}
	return result
}

func fileList(serverName string) []string {
	address, exists := servers[serverName]
	if !exists {
		return []string{"server not registered."}
	}

	conn, err := newConnection(address)
	if err != nil {
		return []string{"could not connect to server: " + err.Error()}
	}
	defer conn.Close()

	c := pb.NewLogTailClient(conn)

	r, err := c.GetFiles(context.Background(), &pb.RequestFile{})
	if err != nil {
		return []string{"could get the file list: " + err.Error()}
	}

	return r.GetName()
}

func logtail(serverName, fileName string) []string {
	address, exists := servers[serverName]
	if !exists {
		return []string{"server not registered."}
	}

	conn, err := newConnection(address)
	if err != nil {
		return []string{"could not connect to server: " + err.Error()}
	}
	defer conn.Close()

	c := pb.NewLogTailClient(conn)

	r, err := c.GetFileContent(context.Background(), &pb.FileName{Name: &fileName})
	if err != nil {
		return []string{"could get the file list: " + err.Error()}
	}

	return r.GetLine()
}

func logTailCommand(command *bot.Cmd) (string, error) {
	var result string

	switch {

	case command.Args[0] == "list":
		result = strings.Join(serverList(), "\n")

	case command.Args[1] == "list":
		result = strings.Join(fileList(command.Args[0]), "\n")

	default:
		result = strings.Join(logtail(command.Args[0], command.Args[1]), "\n")

	}

	return fmt.Sprintf(result), nil
}

func init() {
	bot.RegisterCommand(
		"logtail",
		"Display the tail of a log file",
		"<server|list> <file|list>",
		logTailCommand,
	)
}
