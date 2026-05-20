package grpc

import (
	"context"
	"fmt"
	"time"

	pb "github.com/aakash811/queueflow/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CheckWorkerHealth() {
	conn, err := grpc.Dial(
		"worker-service:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		fmt.Println(
			"gRPC connection failed:",
			err,
		)
		return
	}

	defer conn.Close()

	client := pb.NewWorkerServiceClient(conn)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5 * time.Second,
	)

	defer cancel()

	resp, err := client.HealthCheck(
		ctx,
		&pb.HealthRequest{},
	)

	if err != nil {
		fmt.Println(
			"Health check failed:",
			err,
		)
		return
	}

	fmt.Println(
		"worker health:",
		resp.Status,
	)
}