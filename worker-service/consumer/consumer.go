package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/shared/logger"
	"github.com/aakash811/queueflow/worker-service/circuitbreaker"
	"github.com/aakash811/queueflow/worker-service/kafka"
	"github.com/aakash811/queueflow/worker-service/models"
	"github.com/aakash811/queueflow/worker-service/processor"
	"github.com/aakash811/queueflow/worker-service/repository"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/sony/gobreaker"
	"go.uber.org/zap"
)

func worker(
	workerID int,
	jobs <-chan models.Job,
	wg *sync.WaitGroup,
) {

	fmt.Println("worker started:", workerID)

	for job := range jobs {

		wg.Add(1)

		logger.Log.Info(
			"worker processing job",

			zap.Int("worker_id", workerID),
			zap.String("job_id", job.ID),
		)


		logger.Log.Info(
			"job timeout seconds",
			zap.Int("timeout_seconds", config.AppConfig.JobTimeoutSeconds),
		)

		timeoutCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Duration(
				config.AppConfig.JobTimeoutSeconds,
			)*time.Second,
		)

		_, err := circuitbreaker.JobBreaker.Execute(
			func() (interface{}, error) {

				return nil,
					processor.ProcessJob(
						timeoutCtx,
						job,
					)
			},
		)

		cancel()

		if err == gobreaker.ErrOpenState {

			logger.Log.Info(
				"circuit breaker open, rejecting job",
				zap.String("job_id", job.ID),
			)

			wg.Done()

			continue
		}

		if errors.Is(err, context.DeadlineExceeded) {

			fmt.Println(
				"job timeout exceeded:",
				job.ID,
			)
		}

		if err != nil {

			fmt.Println("processing error:", err)

			job.RetryCount++

			if job.RetryCount <= job.MaxRetries {

				logger.Log.Warn(
					"retrying job",

					zap.String("job_id", job.ID),
					zap.Int("retry_count", job.RetryCount),
				)

				backoff := time.Duration(
					1<<job.RetryCount,
				) * time.Second

				time.Sleep(backoff)

				err = kafka.PublishRetryJobs(job)

				if err != nil {

					fmt.Println(
						"retry publish error:",
						err,
					)
				}

				wg.Done()

				continue
			}

			fmt.Println(
				"max retries exceeded:",
				job.ID,
			)

			err = kafka.PublishDeadLetterJob(job)

			if err != nil {

				fmt.Println(
					"dead letter publish error:",
					err,
				)
			}

			err = repository.SavedDeadLetterJob(job)

			if err != nil {

				fmt.Println(
					"dead letter save error:",
					err,
				)
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