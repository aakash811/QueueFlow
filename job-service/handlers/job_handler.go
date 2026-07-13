package handlers

import (
	"net/http"

	"github.com/aakash811/queueflow/job-service/models"
	"github.com/aakash811/queueflow/job-service/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func CreateJobHandler(c *gin.Context) {
	var job models.Job

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if idemKey, exists := c.Get("idempotency_key"); exists {
		job.IdempotencyKey = idemKey.(string)
	}

	if err := validate.Struct(job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := services.CreateJob(c.Request.Context(), job)

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