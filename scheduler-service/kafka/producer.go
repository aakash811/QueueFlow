package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aakash811/queueflow/scheduler-service/models"
	"github.com/aakash811/queueflow/shared/tracing"
	kafkago "github.com/segmentio/kafka-go"
)

var Producer *kafkago.Writer

func InitProducer() {

	Producer = &kafkago.Writer{
		Addr: kafkago.TCP("kafka:9092"),
		Topic: "jobs_pending",
		Balancer: &kafkago.LeastBytes{},
	}

	fmt.Println("scheduler kafka producer initialized")
}

func PublishJob(ctx context.Context, job models.Job) error {

	payload, err := json.Marshal(job)

	if err != nil {
		return err
	}

	headers := map[string]string{}
	tracing.InjectTraceContext(ctx, headers)
	headers["publish-time"] = time.Now().UTC().Format(time.RFC3339Nano)

	err = Producer.WriteMessages(
		ctx,
		kafkago.Message{
			Value:   payload,
			Headers: convertHeaders(headers),
		},
	)

	return err
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
