package repository

import (
	"context"
	"fmt"

	"github.com/aakash811/queueflow/job-service/db"
	"github.com/aakash811/queueflow/job-service/models"
)

func CreateJob(job models.Job) error {
	query := `
	INSERT INTO jobs (
		id,
		queue_name,
		payload,
		status,
		retry_count,
		max_retries
	)
	VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := db.DB.Exec(
		context.Background(),
		query,
		job.ID,
		job.QueueName,
		job.Payload,
		job.Status,
		job.RetryCount,
		job.MaxRetries,
	)

	if err != nil {
		fmt.Println(err)
	}

	return err
}