package scheduler

import (
	"fmt"
	"time"

	"github.com/aakash811/queueflow/scheduler-service/kafka"
	"github.com/aakash811/queueflow/scheduler-service/repository"
)

func StartScheduler() {

	fmt.Println("scheduler service started")

	ticker := time.NewTicker(5 * time.Second)

	defer ticker.Stop()

	for {

		<-ticker.C

		jobs, err := repository.GetReadyJobs()

		if err != nil {

			fmt.Println(
				"scheduler fetch error:",
				err,
			)

			continue
		}

		for _, job := range jobs {

			fmt.Println(
				"dispatching scheduled job:",
				job.ID,
			)

			err = kafka.PublishJob(job)

			if err != nil {

				fmt.Println(
					"publish error:",
					err,
				)

				continue
			}

			err = repository.MarkJobDispatched(
				job.ID,
			)

			if err != nil {

				fmt.Println(
					"mark dispatched error:",
					err,
				)
			}
		}
	}
}
