package main

import (
	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/shared/logger"
	"go.uber.org/zap"
)

func main() {
	config.LoadConfig()
	logger.InitLogger()

	defer logger.Log.Sync()

	logger.Log.Info("job-service started", 
		zap.String("env", config.AppConfig.AppEnv),
		zap.String("port", config.AppConfig.Port),
	)
}