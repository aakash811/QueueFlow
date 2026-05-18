package routes

import (
	"github.com/aakash811/queueflow/job-service/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.POST("/jobs", handlers.CreateJobHandler)
}