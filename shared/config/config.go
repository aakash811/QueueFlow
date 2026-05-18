package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv        string
	Port          string
	PostgresURL   string
	RedisURL      string
	KafkaBrokers  string
	LogLevel      string
}

var AppConfig Config

func LoadConfig() {
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
		AppEnv:       viper.GetString("APP_ENV"),
		Port:         viper.GetString("PORT"),
		PostgresURL:  viper.GetString("POSTGRES_URL"),
		RedisURL:     viper.GetString("REDIS_URL"),
		KafkaBrokers: viper.GetString("KAFKA_BROKERS"),
		LogLevel:     viper.GetString("LOG_LEVEL"),
	}
}
