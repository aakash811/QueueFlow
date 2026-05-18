package services

import (
	"github.com/aakash811/queueflow/job-service/kafka"
	"github.com/aakash811/queueflow/job-service/models"
	"github.com/aakash811/queueflow/job-service/repository"
	"github.com/google/uuid"
)

func CreateJob(job models.Job) error {
	job.ID = uuid.New().String()
	job.Status = "pending"
	job.RetryCount = 0
	job.MaxRetries = 3

	err := repository.CreateJob(job)

	if err != nil {
		return err
	}

	err = kafka.PublishJob(job)

	if err != nil {
		return err
	}

	return nil
}