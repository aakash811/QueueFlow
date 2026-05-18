package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aakash811/queueflow/shared/config"
	"github.com/aakash811/queueflow/worker-service/models"
	"github.com/aakash811/queueflow/worker-service/processor"
	kafkago "github.com/segmentio/kafka-go"
)

func worker(workerID int, jobs <-chan models.Job) {
	fmt.Println("worker started:", workerID)

	for job := range jobs {
		fmt.Printf(
			"worker %d processing jobs %s\n",
			workerID,
			job.ID,
		)

		err := processor.ProcessJob(job)

		if err != nil {
			fmt.Println("processing error:", err)
		}
	}
}

func StartConsumer() {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic:   "jobs_pending",
		GroupID: "queueflow-workers",
	})

	fmt.Println("worker consumer started")

	jobChannel := make(chan models.Job, 100)

	workerCount := config.AppConfig.WorkerConcurrency

	for i := 1; i <= workerCount; i++ {
		go worker(i, jobChannel)
	}

	for {
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