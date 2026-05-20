package routes

import (
	"github.com/aakash811/queueflow/job-service/handlers"
	"github.com/aakash811/queueflow/job-service/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	router.POST("/login", handlers.LoginHandler)

	router.GET(
		"/dashboard/summary",
		handlers.DashboardSummaryHandler,
	)

	router.GET(
		"/jobs/recent",
		handlers.GetRecentJobsHandler,
	)

	router.GET(
		"/dead-letter/recent",
		handlers.GetDeadLetterJobsHandler,
	)

	authorized := router.Group("/")

	authorized.Use(
		middleware.AuthMiddleware(),
		middleware.RateLimitMiddleware(),
	)

	authorized.POST(
		"/jobs",
		middleware.RequireRole("ADMIN"),
		middleware.BackpressureMiddleware(),
		middleware.IdempotencyMiddleware(),
		handlers.CreateJobHandler,
	)

}