package config

import (
	"log"

	"github.com/aakash811/queueflow/shared/secrets"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv        string
	Port          string
	PostgresURL   string
	RedisURL      string
	KafkaBrokers  string
	LogLevel      string
	WorkerConcurrency int
	JobTimeoutSeconds int
	MaxPendingJobs int
	JWTSecret string
	KafkaDefaultPartitions int
}

var AppConfig Config

func LoadConfig() {
	secrets.LoadSecrets()

	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	AppConfig = Config{
		AppEnv:             viper.GetString("APP_ENV"),
		Port:               viper.GetString("PORT"),
		PostgresURL:        secrets.GetString("POSTGRES_URL", viper.GetString("POSTGRES_URL")),
		RedisURL:           secrets.GetString("REDIS_URL", viper.GetString("REDIS_URL")),
		KafkaBrokers:       secrets.GetString("KAFKA_BROKERS", viper.GetString("KAFKA_BROKERS")),
		LogLevel:           viper.GetString("LOG_LEVEL"),
		WorkerConcurrency:  secrets.GetInt("WORKER_CONCURRENCY", viper.GetInt("WORKER_CONCURRENCY")),
		JobTimeoutSeconds:  secrets.GetInt("JOB_TIMEOUT_SECONDS", viper.GetInt("JOB_TIMEOUT_SECONDS")),
		MaxPendingJobs:     secrets.GetInt("MAX_PENDING_JOBS", viper.GetInt("MAX_PENDING_JOBS")),
		JWTSecret:          secrets.GetString("JWT_SECRET", viper.GetString("JWT_SECRET")),
		KafkaDefaultPartitions: secrets.GetInt("KAFKA_DEFAULT_PARTITIONS", viper.GetInt("KAFKA_DEFAULT_PARTITIONS")),
	}
}
