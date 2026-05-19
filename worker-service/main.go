package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/worker-service/circuitbreaker"
	"github.com/aakash811/queueflow/worker-service/consumer"
	"github.com/aakash811/queueflow/worker-service/db"
	"github.com/aakash811/queueflow/worker-service/kafka"
)

func main() {
	config.LoadConfig()

	err := db.ConnectDatabase(config.AppConfig.PostgresURL)

	if err != nil {
		panic(err)
	}
	kafka.InitProducer()

	ctx, cancel := context.WithCancel(context.Background())

	signalChannel := make(chan os.Signal, 1)

	signal.Notify(
		signalChannel,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		<-signalChannel
		fmt.Println("shutdown initiated")

		cancel()
	}()

	circuitbreaker.InitCircuitBreaker()
	consumer.StartConsumer(ctx)
}