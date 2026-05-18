package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/worker-service/kafka"
	"github.com/aakash811/queueflow/worker-service/models"
	"github.com/aakash811/queueflow/worker-service/processor"
	"github.com/aakash811/queueflow/worker-service/repository"
	kafkago "github.com/segmentio/kafka-go"
)

func worker(
	workerID int,
	jobs <-chan models.Job, 
	wg *sync.WaitGroup,
) {
	fmt.Println("worker started:", workerID)

	for job := range jobs {
		fmt.Printf(
			"worker %d processing job %s\n",
			workerID,
			job.ID,
		)

		wg.Add(1)
		err := processor.ProcessJob(job)

		if err != nil {
			fmt.Println("processing error:", err)

			job.RetryCount++

			if job.RetryCount <= job.MaxRetries {
				fmt.Printf(
					"retrying job %s attempt %d\n",
					job.ID,
					job.RetryCount,
				)

				backoff := time.Duration(
					1 << job.RetryCount,
				) * time.Second

				time.Sleep(backoff)

				err = kafka.PublishRetryJobs(job)

				if err != nil {
					fmt.Println("retry publish error:", err)
				}

				continue
			}
			fmt.Println("max retries exceeded:", job.ID)

			err = kafka.PublishDeadLetterJob(job)

			if err != nil {
				fmt.Println("dead letter publish error:", err)
			}

			err = repository.SavedDeadLetterJob(job)

			if err != nil {
				fmt.Println("dead letter save error:", err)
			}
		}

		wg.Done()
	}
}

func StartConsumer(ctx context.Context) {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "jobs_pending",
		GroupID: "queueflow-workers",
	})

	fmt.Println("worker consumer started")
	
	jobChannel := make(chan models.Job, 100)
	var wg sync.WaitGroup

	workerCount := config.AppConfig.WorkerConcurrency

	for i := 1; i <= workerCount; i++ {
		go worker(i, jobChannel, &wg)
	}

	go StartRetryConsumer(jobChannel)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("shutdown signal received, shutting down consumer...")
			close(jobChannel)

			wg.Wait()
			err := reader.Close()

			if err != nil {
				fmt.Println("Reader close error:", err)
			}

			fmt.Println("worker shutdown complete")

			return
		default:
		}
		message, err := reader.ReadMessage(context.Background())

		if err != nil {
			fmt.Println("Consumer error:", err)
			continue
		}

		var job models.Job
		err = json.Unmarshal(message.Value, &job)

		if err != nil {
			fmt.Println("json parse error:", err)
			continue
		}

		jobChannel <- job
	}
}