package cronjobs

import (
	"fmt"

	"github.com/aakash811/queueflow/scheduler-service/kafka"
	"github.com/aakash811/queueflow/scheduler-service/models"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

func StartCronJobs() {
	c := cron.New()

	_, err := c.AddFunc(
		"@every 1m",
		func ()  {
			job := models.Job{
				ID: uuid.New().String(),
				QueueName: "cron-email",
				Status: "pending",
				RetryCount: 0,
				MaxRetries: 3,
				Payload: map[string]interface{}{
					"type": "recurring-cron-job",
				},
			}

			fmt.Println("running recurring cron job:", job.ID)

			err := kafka.PublishJob(job)

			if err != nil {
				fmt.Println(

					"cron publish error:",
					err,
				)
				return
			}

			fmt.Println(
				"cron job dispatched:",
				job.ID,
			)
		},
	)

	if err != nil {
		fmt.Println(
			"cron registration error:",
			err,
		)
		return
	}

	c.Start()

	fmt.Println("Cron Scheduler Started")
}