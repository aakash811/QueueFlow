package services

import (
	"github.com/aakash811/queueflow/job-service/models"
	"github.com/aakash811/queueflow/job-service/repository"
)

func CreateJob(job models.Job) error {
	job.Status = "pending"
	job.RetryCount = 0
	job.MaxRetries = 3

	return repository.CreateJob(job)
}