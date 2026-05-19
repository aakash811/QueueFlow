package main

import (
	"fmt"

	"github.com/aakash811/queueflow/shared/config"

	cronjobs "github.com/aakash811/queueflow/scheduler-service/cron"
	"github.com/aakash811/queueflow/scheduler-service/db"
	"github.com/aakash811/queueflow/scheduler-service/kafka"
	"github.com/aakash811/queueflow/scheduler-service/scheduler"
)

func main() {

	config.LoadConfig()

	err := db.ConnectDatabase(
		config.AppConfig.PostgresURL,
	)

	if err != nil {
		panic(err)
	}

	kafka.InitProducer()

	fmt.Println("scheduler initialized")

	go cronjobs.StartCronJobs()
	scheduler.StartScheduler()
}
