package main

import (
	"context"
	"github.com/aakash811/queueflow/job-service/db"
	"github.com/aakash811/queueflow/job-service/grpc"
	"github.com/aakash811/queueflow/job-service/kafka"
	"github.com/aakash811/queueflow/job-service/redis"
	"github.com/aakash811/queueflow/job-service/routes"
	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/shared/logger"
	"github.com/aakash811/queueflow/shared/tracing"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

func main() {
	err := logger.InitLogger()
	if err != nil {
		panic(err)
	}

	defer logger.Log.Sync()

	config.LoadConfig()

	shutdown := tracing.InitTracer("job-service")
	defer shutdown(context.Background())

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
	router.Use(otelgin.Middleware("job-service"))
	routes.RegisterRoutes(router)
	router.GET("/metrics", func(c *gin.Context) {
		promhttp.Handler().ServeHTTP(c.Writer, c.Request)
	})
	logger.Log.Info(
		"job-service started",
		zap.String(
			"port",
			config.AppConfig.Port,
		),
	)
	router.Run(":" + config.AppConfig.Port)
}