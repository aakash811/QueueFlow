package kafka

import (
	"context"
	"encoding/json"

	"github.com/aakash811/queueflow/job-service/models"
	kafkago "github.com/segmentio/kafka-go"
)

var Writer *kafkago.Writer

func InitProducer() {
	Writer = &kafkago.Writer{
		Addr:     kafkago.TCP("kafka:9092"),
		Topic:    "jobs_pending",
		Balancer: &kafkago.LeastBytes{},
	}
}

func PublishJob(job models.Job) error {
	payload, err := json.Marshal(job)

	if err != nil {
		return err
	}

	return Writer.WriteMessages(
		context.Background(),
		kafkago.Message{
			Key:   []byte(job.ID),
			Value: payload,
		},
	)
}