package services

import (
	"fmt"

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

	if job.ExecuteAt != nil {
		fmt.Println(
			"Scheduled job stored for later execution:",
			job.ID,
		)

		return nil
	}

	fmt.Println("Publishing immediate job:", job.ID)
	err = kafka.PublishJob(job)

	if err != nil {
		return err
	}

	return nil
}