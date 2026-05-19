package middleware

import (
	"net/http"

	"github.com/aakash811/queueflow/job-service/repository"
	"github.com/aakash811/queueflow/shared/config"
	"github.com/gin-gonic/gin"
)

func BackpressureMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		count, err := repository.CountPendingJobs()

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "queue pressure check failed",
			})

			ctx.Abort()
			return
		}

		if count >= config.AppConfig.MaxPendingJobs {
			ctx.JSON(http.StatusTooManyRequests, gin.H{
				"error": "queue overloaded",
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}