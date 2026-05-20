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

func GetRecentJobs() ([]models.Job, error) {

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
		execute_at,
		processed_at,
		failed_at
	FROM jobs
	ORDER BY created_at DESC
	LIMIT 20
	`

	rows, err := db.DB.Query(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	jobs := []models.Job{}

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

func GetDashboardSummary() (map[string]interface{}, error) {

	var queueDepth int
	var completed int
	var failed int
	var retries int

	err := db.DB.QueryRow(
		context.Background(),
		`
		SELECT COUNT(*)
		FROM jobs
		WHERE status = 'pending'
		`,
	).Scan(&queueDepth)

	if err != nil {
		return nil, err
	}

	err = db.DB.QueryRow(
		context.Background(),
		`
		SELECT COUNT(*)
		FROM jobs
		WHERE status = 'completed'
		`,
	).Scan(&completed)

	if err != nil {
		return nil, err
	}

	err = db.DB.QueryRow(
		context.Background(),
		`
		SELECT COUNT(*)
		FROM jobs
		WHERE status = 'failed'
		`,
	).Scan(&failed)

	if err != nil {
		return nil, err
	}

	err = db.DB.QueryRow(
		context.Background(),
		`
		SELECT COALESCE(SUM(retry_count), 0)
		FROM jobs
		`,
	).Scan(&retries)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"queue_depth": queueDepth,
		"throughput": completed,
		"retries": retries,
		"failures": failed,
		"latency": 120,
		"workers": 4,
	}, nil
}

func GetDeadLetterJobs() ([]models.DeadLetterJob, error) {

	query := `
	SELECT *
	FROM dead_letter_jobs
	ORDER BY failed_at DESC
	LIMIT 20
	`

	rows, err := db.DB.Query(
		context.Background(),
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	jobs := []models.DeadLetterJob{}

	for rows.Next() {

		job := models.DeadLetterJob{}

		err := rows.Scan(
			&job.ID,
			&job.QueueName,
			&job.Payload,
			&job.RetryCount,
			&job.FailedAt,
		)

		if err != nil {
			return nil, err
		}

		jobs = append(jobs, job)
	}

	return jobs, nil
}