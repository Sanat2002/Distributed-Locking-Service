package main

import (
	"context"
	"log"
	"time"

	pb "github.com/username/distributed-lock-service/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewLockServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.Ping(ctx, &pb.PingRequest{
		Message: "Hello from client",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Response:", resp.Message)
}
