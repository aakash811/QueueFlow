package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/aakash811/queueflow/job-service/redis"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()

		key := "rate_limit" + clientIP

		ctx := context.Background()

		count, err := redis.Client.Incr(ctx, key).Result()

		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": "rate limiter failed",
				},
			)
			c.Abort()
			return
		}

		if count == 1 {
			redis.Client.Expire(ctx, key, time.Minute)
		}

		if count > 5 {
			c.JSON(
				http.StatusTooManyRequests,
				gin.H{
					"error": "rate limit exceeded",
				},
			)
			c.Abort()
			return		
		}

		c.Next()
	}
}