package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aakash811/queueflow/worker-service/models"
	kafkago "github.com/segmentio/kafka-go"
)

func StartRetryConsumer(jobChannel chan models.Job) {
	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{"kafka:9092"},
		Topic: "jobs_retry",
		GroupID: "queueflow-retry-workers",
	})

	fmt.Println("retry consumer started")

	for {
		message, err := reader.ReadMessage(context.Background())

		if err != nil {
			fmt.Println("retry consumer error:", err)
			continue
		}

		var job models.Job

		err = json.Unmarshal(message.Value, &job)

		if err != nil {
			fmt.Println("retry json parse error:", err)
			continue
		}

		jobChannel <- job
	}
}