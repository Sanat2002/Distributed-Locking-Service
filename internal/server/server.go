package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc/reflection"

	pb "github.com/username/distributed-lock-service/internal/proto"

	"github.com/username/distributed-lock-service/internal/lock"
	rds "github.com/username/distributed-lock-service/internal/redis"

	"google.golang.org/grpc"
)


type server struct {
	pb.UnimplementedLockServiceServer
	lockManager *lock.Manager
}

func (s *server) Ping(
	ctx context.Context,
	req *pb.PingRequest,
) (*pb.PingResponse, error) {

	// 1. Check Redis connectivity
	if err := s.lockManager.HealthCheck(ctx); err != nil {
		return nil, err
	}
	// 2. If we reached here:
	// - gRPC is working
	// - Redis is reachable
	return &pb.PingResponse{
		Message: "client->server OK, server->redis OK",
	}, nil
}


func main() {
	redisClient := rds.NewClient("localhost:6379")
	log.Println("connected to redis")

	lockManager := lock.NewManager(redisClient)

	srv := &server{
		lockManager: lockManager,
	}
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterLockServiceServer(grpcServer, srv)

	reflection.Register(grpcServer)
	
	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
