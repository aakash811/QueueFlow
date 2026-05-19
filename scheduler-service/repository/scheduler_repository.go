package repository

import (
	"context"

	"github.com/aakash811/queueflow/scheduler-service/db"
	"github.com/aakash811/queueflow/scheduler-service/models"
)

func GetReadyJobs() ([]models.Job, error) {

	query := `
	SELECT
		id,
		queue_name,
		payload,
		status,
		retry_count,
		max_retries,
		created_at,
		updated_at,
		scheduled_at,
		processed_at,
		failed_at
	FROM jobs
	WHERE scheduled_at IS NOT NULL
	AND scheduled_at <= NOW()
	AND status = 'pending'
	`

	rows, err := db.DB.Query(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var jobs []models.Job

	for rows.Next() {

		var job models.Job

		err := rows.Scan(
			&job.ID,
			&job.QueueName,
			&job.Payload,
			&job.Status,
			&job.RetryCount,
			&job.MaxRetries,
			&job.CreatedAt,
			&job.UpdatedAt,
			&job.ExecuteAt,
			&job.ProcessedAt,
			&job.FailedAt,
		)

		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}

func MarkJobDispatched(jobID string) error {

	query := `
	UPDATE jobs
	SET scheduled_at = NULL
	WHERE id = $1
	`

	_, err := db.DB.Exec(
		context.Background(),
		query,
		jobID,
	)

	return err
}
