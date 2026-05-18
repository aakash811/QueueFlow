package main

import (
	"fmt"

	"github.com/aakash811/queueflow/job-service/db"
	"github.com/aakash811/queueflow/job-service/kafka"
	"github.com/aakash811/queueflow/job-service/routes"
	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/shared/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	logger.InitLogger()

	defer logger.Log.Sync()

	err := db.ConnectDatabase(config.AppConfig.PostgresURL)

	if err != nil {
		panic(err)
	}

	kafka.InitProducer() 
	router := gin.Default()
	routes.RegisterRoutes(router)
	logger.Log.Info("job-service started")
	fmt.Println(config.AppConfig.PostgresURL)
	router.Run(":" + config.AppConfig.Port)
}