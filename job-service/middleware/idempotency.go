package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/aakash811/queueflow/job-service/redis"
	"github.com/gin-gonic/gin"
)

func IdempotencyMiddleware() gin.HandlerFunc {
	return func (ctx *gin.Context) {
		key := ctx.GetHeader("Idempotency-Key")
		
		if key == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "missing idempotency key",
			})
			ctx.Abort()

			return
		}

		lockKey := "idem:" + key

		success, err := redis.Client.SetNX(
			context.Background(),
			lockKey,
			"locked",
			60 * time.Second,
		).Result()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Idempotency check failed",
			})
			ctx.Abort()
			return
		}

		if !success {
			ctx.JSON(http.StatusConflict, gin.H{
				"error": "duplicate request",
			})
			ctx.Abort()
			return
		}

		ctx.Set("idempotency_key", key)
		ctx.Next()
	}
}