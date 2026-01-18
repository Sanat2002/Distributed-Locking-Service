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

	client := pb.NewReadwriteservicesClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var i int64

	for i = 0; i < 5; i++ {
		resp, err := client.Read(ctx, &pb.ReadRequest{
			Read: true,
		})
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Read Response:", resp.Result, resp.CurrData)

		resp1, err1 := client.Write(ctx, &pb.WriteRequest{
			ResourceId: "user:1",
			Val: 100 * i,
		})

		if err1 != nil {
			log.Fatal(err1)
		}

		log.Println("Write Response:", resp1.Result, resp1.UpdatedData)
	}
}
