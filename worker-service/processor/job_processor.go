package processor

import (
	"fmt"
	"time"

	"github.com/aakash811/queueflow/worker-service/models"
	"github.com/aakash811/queueflow/worker-service/repository"
)

func ProcessJob(job models.Job) error {
	fmt.Println("Processing job:", job.ID)

	time.Sleep(2 * time.Second)

	err := repository.UpdateJobStatus(job.ID, "completed")

	if err != nil {
		return err
	}

	fmt.Println("Job completed:", job.ID)

	return nil
}