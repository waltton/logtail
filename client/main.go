package main

import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/waltton/logtail/logtail"
)

const address = "localhost:50051"

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewLogTailClient(conn)

	// Contact the server and print out its response.
	r, err := c.GetFiles(context.Background(), &pb.RequestFile{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Print(r.GetName())
}
