package main

import (
	"context"
	"net/http"

	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/shared/logger"
	"github.com/aakash811/queueflow/shared/metrics"
	"github.com/aakash811/queueflow/shared/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	cronjobs "github.com/aakash811/queueflow/scheduler-service/cron"
	"github.com/aakash811/queueflow/scheduler-service/db"
	"github.com/aakash811/queueflow/scheduler-service/kafka"
	"github.com/aakash811/queueflow/scheduler-service/scheduler"
)

func main() {

	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

	config.LoadConfig()

	shutdown := tracing.InitTracer("scheduler-service")
	defer shutdown(context.Background())

	err = db.ConnectDatabase(
		config.AppConfig.PostgresURL,
	)

	if err != nil {
		panic(err)
	}

	kafka.InitProducer()

	logger.Log.Info(
		"scheduler initialized",
		zap.String(
			"port",
			config.AppConfig.Port,
		),
	)

	go cronjobs.StartCronJobs()

	metrics.InitMetrics()
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(":2113", nil)

	scheduler.StartScheduler()
}
