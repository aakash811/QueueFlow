package grpc

import (
	"context"
	"fmt"
	"net"

	pb "github.com/aakash811/queueflow/proto"
	"github.com/aakash811/queueflow/worker-service/repository"
	"google.golang.org/grpc"
)

type WorkerServer struct {
	pb.UnimplementedWorkerServiceServer
}

func (s *WorkerServer) UpdateJobStatus(
	ctx context.Context,
	req *pb.JobStatusRequest,
) (*pb.JobStatusResponse, error) {
	
	err := repository.UpdateJobStatus(
		req.JobId,
		req.Status,
	)

	if err != nil {
		return nil, err
	}

	return &pb.JobStatusResponse{
		Message: "Job status updated",
	}, nil
}

func (s *WorkerServer) HealthCheck(
	ctx context.Context,
	req *pb.HealthRequest,
) (*pb.HealthResponse, error) {
	
	return &pb.HealthResponse{
		Status: "Healthy",
	}, nil
}

func StartGRPCServer() {
	lis, err := net.Listen(
		"tcp",
		":50051",
	)

	if err != nil {
		panic(err)
	}


	grpcServer := grpc.NewServer()

	pb.RegisterWorkerServiceServer(
		grpcServer,
		&WorkerServer{},
	)

	fmt.Println(
		"gRPC server started on port 50051",
	)

	err = grpcServer.Serve(lis)

	if err != nil {
		panic(err)
	}
}