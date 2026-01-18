package main

import (
	"context"
	"log"
	"net"

	pb "github.com/username/distributed-lock-service/internal/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedLockServiceServer
}

func (s *server) Ping(
	ctx context.Context,
	req *pb.PingRequest,
) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Message: "Server received: " + req.Message,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLockServiceServer(grpcServer, &server{})

	log.Println("gRPC server listening on :50051")
	grpcServer.Serve(lis)
}
