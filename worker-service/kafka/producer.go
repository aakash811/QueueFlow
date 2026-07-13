package kafka

import (
	"context"
	"encoding/json"

	"github.com/aakash811/queueflow/shared/tracing"
	"github.com/aakash811/queueflow/worker-service/models"
	kafkago "github.com/segmentio/kafka-go"
)

var RetryWriter *kafkago.Writer
var DeadLetterWriter *kafkago.Writer

func InitProducer() {
	RetryWriter = &kafkago.Writer{
		Addr: kafkago.TCP("kafka:9092"),
		Topic: "jobs_retry",
		Balancer: &kafkago.LeastBytes{},
	}

	DeadLetterWriter = &kafkago.Writer{
		Addr: kafkago.TCP("kafka:9092"),
		Topic: "jobs_deadletter",
		Balancer: &kafkago.LeastBytes{},
	}
}

func PublishRetryJobs(ctx context.Context, job models.Job) error {
	payload, err := json.Marshal(job)

	if err != nil {
		return err
	}

	headers := map[string]string{}
	tracing.InjectTraceContext(ctx, headers)

	return RetryWriter.WriteMessages(
		ctx,
		kafkago.Message{
			Key:     []byte(job.QueueName),
			Value:   payload,
			Headers: convertHeaders(headers),
		},
	)
}

func PublishDeadLetterJob(ctx context.Context, job models.Job) error {
	payload, err := json.Marshal(job)

	if err != nil {
		return err
	}

	headers := map[string]string{}
	tracing.InjectTraceContext(ctx, headers)

	return DeadLetterWriter.WriteMessages(
		ctx,
		kafkago.Message{
			Key:     []byte(job.QueueName),
			Value:   payload,
			Headers: convertHeaders(headers),
		},
	)
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