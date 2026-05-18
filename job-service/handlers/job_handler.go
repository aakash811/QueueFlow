package handlers

import (
	"net/http"

	"github.com/aakash811/queueflow/job-service/models"
	"github.com/aakash811/queueflow/job-service/services"
	"github.com/gin-gonic/gin"
)

func CreateJobHandler(c *gin.Context) {
	var job models.Job

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := services.CreateJob(job)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Job Created Successfully",
	})
}