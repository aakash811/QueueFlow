package processor

import (
	"errors"
	"fmt"
	"time"

	"github.com/aakash811/queueflow/worker-service/models"
	"github.com/aakash811/queueflow/worker-service/repository"
)

func ProcessJob(job models.Job) error {
	fmt.Println("Processing job:", job.ID)

	time.Sleep(2 * time.Second)

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