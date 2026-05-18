package repository

import (
	"context"
	"time"

	"github.com/aakash811/queueflow/worker-service/db"
)

func UpdateJobStatus(jobId string, status string) error {
	query := `
	UPDATE jobs
	SET
		status = $1,
		updated_at = NOW(),
		processed_at = $2
	WHERE id = $3
	`

	_, err := db.DB.Exec(
		context.Background(),
		query,
		status,
		time.Now(),
		jobId,
	)

	return err
}