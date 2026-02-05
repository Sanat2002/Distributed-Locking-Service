package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/reflection"

	pb "github.com/username/distributed-lock-service/internal/proto"

	"github.com/username/distributed-lock-service/internal/lock"
	rds "github.com/username/distributed-lock-service/internal/redis"

	"google.golang.org/grpc"
)


type server struct {
	lockManager *lock.Manager
	pb.ReadwriteservicesClient
	pb.ReadwriteservicesServer
}

func (s *server) Read(
	ctx context.Context,
	req *pb.ReadRequest,
) (*pb.ReadResponse, error) {
  	// 1. Check Redis connectivity
	if err := s.lockManager.HealthCheck(ctx); err != nil {
		return nil, err
	}
	// 2. If we reached here:
	// - gRPC is working
	// - Redis is reachable

	return &pb.ReadResponse{
		Result:   "OK",
		CurrData: 200,
	}, nil
}

func (s *server) Write(
	ctx context.Context,
	req *pb.WriteRequest,
) (*pb.WriteResponse, error) {

	log.Println("[WRITE] request received",
		"resource:", req.ResourceId,
		"value:", req.Val,
	)
  
	token, err := s.lockManager.Acquire(
		ctx,
		req.ResourceId,
		5*time.Second,
	)
	if err != nil {
		log.Println("[WRITE] lock acquisition failed:", err)
		return &pb.WriteResponse{
			Result: "LOCK_BUSY",
		}, nil
	}

	log.Println("[WRITE] lock acquired",
		"resource:", req.ResourceId,
		"token:", token,
	)

	updated := 100 + req.Val
	log.Println("[WRITE] write executed",
		"resource:", req.ResourceId,
		"updated_value:", updated,
	)

	// STEP 3: Return success
	return &pb.WriteResponse{
		Result:      "OK",
		UpdatedData: updated,
	}, nil

	// // 1. Check Redis connectivity
	// if err := s.lockManager.HealthCheck(ctx); err != nil {
	// 	return nil, err
	// }
	// // 2. If we reached here:
	// // - gRPC is working
	// // - Redis is reachable
	
	// return &pb.WriteResponse{
	// 	Result:      "OK",
	// 	UpdatedData: 100 + req.Val,
	// 	}, nil
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

	pb.RegisterReadwriteservicesServer(grpcServer, srv)

	reflection.Register(grpcServer)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
