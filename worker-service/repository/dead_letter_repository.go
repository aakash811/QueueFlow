package repository

import (
	"context"
	"encoding/json"

	"github.com/aakash811/queueflow/worker-service/db"
	"github.com/aakash811/queueflow/worker-service/models"
	"github.com/google/uuid"
)

func SavedDeadLetterJob(job models.Job) error {
	payload, err := json.Marshal(job.Payload)

	if err != nil {
		return err
	}

	query := `
		INSERT INTO dead_letter_jobs (
			id,
			original_job_id,
			payload,
			failure_reason,
			failed_at
		)
		VALUES ($1, $2, $3, $4, NOW())
	`

	_, err = db.DB.Exec(
		context.Background(),
		query,
		uuid.New(),
		job.ID,
		payload,
		"max retries exceeded",
	)

	return err
}