package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aakash811/queueflow/shared/logger"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"go.uber.org/zap"
)

type SecretValue struct {
	PostgresURL   string `json:"POSTGRES_URL"`
	RedisURL      string `json:"REDIS_URL"`
	KafkaBrokers  string `json:"KAFKA_BROKERS"`
	JWTSecret     string `json:"JWT_SECRET"`
	DBPassword    string `json:"DB_PASSWORD"`
}

var FetchedSecret *SecretValue

func LoadSecrets() {
	secretID := os.Getenv("AWS_SECRETS_MANAGER_SECRET_ID")

	if secretID == "" {
		logger.Log.Info("AWS_SECRETS_MANAGER_SECRET_ID not set, using environment variables")
		FetchedSecret = nil
		return
	}

	ctx := context.Background()

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.Log.Fatal("failed to load AWS config", zap.Error(err))
	}

	client := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}

	result, err := client.GetSecretValue(ctx, input)
	if err != nil {
		logger.Log.Fatal("failed to fetch secret from AWS Secrets Manager", zap.Error(err))
	}

	var secret SecretValue
	if err := json.Unmarshal([]byte(*result.SecretString), &secret); err != nil {
		logger.Log.Fatal("failed to parse secret JSON", zap.Error(err))
	}

	FetchedSecret = &secret
	logger.Log.Info("loaded secrets from AWS Secrets Manager", zap.String("secret_id", secretID))
}

func GetString(key, fallback string) string {
	if FetchedSecret != nil {
		switch key {
		case "POSTGRES_URL":
			if FetchedSecret.PostgresURL != "" {
				return FetchedSecret.PostgresURL
			}
		case "REDIS_URL":
			if FetchedSecret.RedisURL != "" {
				return FetchedSecret.RedisURL
			}
		case "KAFKA_BROKERS":
			if FetchedSecret.KafkaBrokers != "" {
				return FetchedSecret.KafkaBrokers
			}
		case "JWT_SECRET":
			if FetchedSecret.JWTSecret != "" {
				return FetchedSecret.JWTSecret
			}
		case "DB_PASSWORD":
			if FetchedSecret.DBPassword != "" {
				return FetchedSecret.DBPassword
			}
		}
	}

	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func GetInt(key string, fallback int) int {
	val := GetString(key, "")
	if val == "" {
		return fallback
	}

	var result int
	if _, err := fmt.Sscanf(val, "%d", &result); err != nil {
		log.Printf("invalid int for %s: %v", key, err)
		return fallback
	}

	return result
}

func MustGetenv(key string) string {
	val := GetString(key, "")
	if val == "" {
		logger.Log.Fatal("missing required secret", zap.String("key", key))
	}
	return val
}
