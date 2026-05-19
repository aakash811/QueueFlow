package processor

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aakash811/queueflow/worker-service/models"
	"github.com/aakash811/queueflow/worker-service/repository"
)

func ProcessJob(ctx context.Context, job models.Job) error {
	fmt.Println("Processing job:", job.ID)

	duration := 5
	if value, ok := job.Payload["duration"].(float64); ok {
		duration = int(value)
	}
	fmt.Println("job duration:", duration)
	select {
	case <-time.After(time.Duration(duration) * time.Second):
	case <-ctx.Done():
		fmt.Println("job cancelled:", job.ID)
		return ctx.Err()
	}

	fail, ok := job.Payload["fail"].(bool)

	if ok && fail{
		fmt.Println("Job Failed:", job.ID)
		return errors.New("simulated job failure")
	}

	err := repository.UpdateJobStatus(job.ID, "completed")

	if err != nil {
		return err
	}

	fmt.Println("Job completed:", job.ID)

	return nil
}