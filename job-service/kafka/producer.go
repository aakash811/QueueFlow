package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aakash811/queueflow/job-service/models"
	"github.com/aakash811/queueflow/shared/tracing"
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

func PublishJob(ctx context.Context, job models.Job) error {
	payload, err := json.Marshal(job)

	if err != nil {
		return err
	}

	headers := map[string]string{}
	tracing.InjectTraceContext(ctx, headers)
	headers["publish-time"] = time.Now().UTC().Format(time.RFC3339Nano)

	msg := kafkago.Message{
		Key:     []byte(job.QueueName),
		Value:   payload,
		Headers: convertHeaders(headers),
	}

	return Writer.WriteMessages(ctx, msg)
}

func convertHeaders(headers map[string]string) []kafkago.Header {
	var kafkaHeaders []kafkago.Header
	for k, v := range headers {
		kafkaHeaders = append(kafkaHeaders, kafkago.Header{
			Key:   k,
			Value: []byte(v),
		})
	}
	return kafkaHeaders
}