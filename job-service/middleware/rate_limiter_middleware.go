package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/aakash811/queueflow/job-service/db"
	"github.com/aakash811/queueflow/job-service/redis"
	"github.com/aakash811/queueflow/shared/config"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware() gin.HandlerFunc {
	limit := config.AppConfig.MaxPendingJobs
	if limit <= 0 {
		limit = 100
	}

	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")
		rateLimitKey := ""

		if apiKey != "" {
			hash := sha256.Sum256([]byte(apiKey))
			keyHash := hex.EncodeToString(hash[:])

			var tenantID string
			err := db.DB.QueryRow(
				c.Request.Context(),
				"SELECT tenant_id FROM api_keys WHERE key_hash = $1 AND is_active = TRUE",
				keyHash,
			).Scan(&tenantID)

			if err == nil {
				rateLimitKey = "rate_limit:tenant:" + tenantID
			}
		}

		if rateLimitKey == "" {
			clientIP := c.ClientIP()
			rateLimitKey = "rate_limit:ip:" + clientIP
		}

		ctx := context.Background()

		count, err := redis.Client.Incr(ctx, rateLimitKey).Result()

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
			redis.Client.Expire(ctx, rateLimitKey, time.Minute)
		}

		if count > int64(limit) {
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