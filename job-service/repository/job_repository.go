package repository

import (
	"context"
	"fmt"

	"github.com/aakash811/queueflow/job-service/db"
	"github.com/aakash811/queueflow/job-service/models"
)

func CountPendingJobs() (int, error) {
	query := `
	SELECT COUNT(*)
	FROM jobs
	WHERE status = 'pending'
	`

	var count int
	
	err := db.DB.QueryRow(
		context.Background(),
		query,
	).Scan(&count)

	return count, err
}


func CreateJob(job models.Job) error {
	query := `
	INSERT INTO jobs (
		id,
		queue_name,
		payload,
		status,
		retry_count,
		max_retries,
		scheduled_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
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
		job.ExecuteAt,
	)

	if err != nil {
		fmt.Println(err)
	}

	return err
}