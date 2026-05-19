package main

import (
	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/shared/logger"
	"go.uber.org/zap"

	cronjobs "github.com/aakash811/queueflow/scheduler-service/cron"
	"github.com/aakash811/queueflow/scheduler-service/db"
	"github.com/aakash811/queueflow/scheduler-service/kafka"
	"github.com/aakash811/queueflow/scheduler-service/scheduler"
)

func main() {

	config.LoadConfig()
	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

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
	scheduler.StartScheduler()
}
