package kafka

import (
	"context"
	"encoding/json"

	"github.com/aakash811/queueflow/worker-service/models"
	kafkago "github.com/segmentio/kafka-go"
)

var RetryWriter *kafkago.Writer

func InitProducer() {
	RetryWriter = &kafkago.Writer{
		Addr: kafkago.TCP("kafka:9092"),
		Topic: "jobs_retry",
		Balancer: &kafkago.LeastBytes{},
	}
}

func PublishRetryJobs(job models.Job) error {
	payload, err := json.Marshal(job)

	if err != nil {
		return err
	}

	return RetryWriter.WriteMessages(
		context.Background(),
		kafkago.Message{
			Key:   []byte(job.ID),
			Value: payload,
		},
	)
}