package handlers

import (
	"net/http"

	"github.com/aakash811/queueflow/job-service/repository"
	"github.com/gin-gonic/gin"
)

func DashboardSummaryHandler(c *gin.Context) {

	data, err := repository.GetDashboardSummary()

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		data,
	)
}

func GetRecentJobsHandler(c *gin.Context) {

	jobs, err := repository.GetRecentJobs()

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		jobs,
	)
}

func GetDeadLetterJobsHandler(c *gin.Context) {

	jobs, err := repository.GetDeadLetterJobs()

	if err != nil {

		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		jobs,
	)
}