package main

import (
	"github.com/aakash811/queueflow/job-service/db"
	"github.com/aakash811/queueflow/job-service/grpc"
	"github.com/aakash811/queueflow/job-service/kafka"
	"github.com/aakash811/queueflow/job-service/redis"
	"github.com/aakash811/queueflow/job-service/routes"
	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/shared/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()
	err := logger.InitLogger()
	
	if err != nil {
		panic(err)
	}

	defer logger.Log.Sync()

	err = db.ConnectDatabase(config.AppConfig.PostgresURL)

	if err != nil {
		panic(err)
	}

	err = redis.ConnectRedis( config.AppConfig.RedisURL)

	if err != nil {
		panic(err)
	}
	
	kafka.InitProducer() 
	go grpc.CheckWorkerHealth()
	
	router := gin.Default()
	router.Use(cors.Default())
	routes.RegisterRoutes(router)
	logger.Log.Info(
		"job-service started",

		zap.String(
			"port",
			config.AppConfig.Port,
		),
	)
	router.Run(":" + config.AppConfig.Port)
}