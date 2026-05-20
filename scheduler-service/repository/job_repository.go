package repository

import (
	"context"

	"github.com/aakash811/queueflow/scheduler-service/db"
	"github.com/aakash811/queueflow/scheduler-service/models"
)

func CreateJob(job models.Job) error {

	query := `
	INSERT INTO jobs (
		id,
		queue_name,
		payload,
		status,
		retry_count,
		max_retries,
		created_at,
		updated_at
	)
	VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		NOW(),
		NOW()
	)
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

	return err
}
