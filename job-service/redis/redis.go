package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

var Client *goredis.Client

func ConnectRedis(redisURL string) error {
	client := goredis.NewClient(&goredis.Options{
		Addr: redisURL,
	})

	_, err := client.Ping(context.Background()).Result()

	if err != nil {
		return err
	}

	fmt.Println("redis connected")

	Client = client
	return nil
}