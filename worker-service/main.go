package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/shared/logger"
	"github.com/aakash811/queueflow/shared/metrics"
	"github.com/aakash811/queueflow/shared/tracing"
	"github.com/aakash811/queueflow/worker-service/circuitbreaker"
	"github.com/aakash811/queueflow/worker-service/consumer"
	"github.com/aakash811/queueflow/worker-service/db"
	"github.com/aakash811/queueflow/worker-service/grpc"
	"github.com/aakash811/queueflow/worker-service/kafka"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

	config.LoadConfig()

	shutdown := tracing.InitTracer("worker-service")
	defer shutdown(context.Background())

	go grpc.StartGRPCServer()
	metrics.InitMetrics()

	http.Handle(
		"/metrics",
		promhttp.Handler(),
	)

	go http.ListenAndServe(
		":2112",
		nil,
	)

	if err != nil {
		panic(err)
	}
	err = db.ConnectDatabase(config.AppConfig.PostgresURL)

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
		logger.Log.Info("shutdown initiated")

		cancel()
	}()

	circuitbreaker.InitCircuitBreaker()
	consumer.StartConsumer(ctx)
}