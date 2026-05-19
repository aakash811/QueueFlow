package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aakash811/queueflow/scheduler-service/models"

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

func PublishJob(job models.Job) error {

	payload, err := json.Marshal(job)

	if err != nil {
		return err
	}

	err = Producer.WriteMessages(
		context.Background(),
		kafkago.Message{
			Value: payload,
		},
	)

	return err
}
