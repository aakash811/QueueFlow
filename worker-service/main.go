package main

import (
	"github.com/aakash811/queueflow/shared/config"
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
	consumer.StartConsumer()
}