package main

import (
	"context"
	"log"
	"net"

	pb "github.com/username/distributed-lock-service/internal/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.ReadwriteservicesClient
	pb.ReadwriteservicesServer
}

func (s *server) Read(
	ctx context.Context,
	req *pb.ReadRequest,
) (*pb.ReadResponse, error) {

	return &pb.ReadResponse{
		Result:   "OK",
		CurrData: 200,
	}, nil
}

func (s *server) Write(
	ctx context.Context,
	req *pb.WriteRequest,
) (*pb.WriteResponse, error) {

	return &pb.WriteResponse{
		Result:      "OK",
		UpdatedData: 100 + req.Val,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterReadwriteservicesServer(grpcServer, &server{})

	log.Println("gRPC server listening on :50051")
	grpcServer.Serve(lis)
}
