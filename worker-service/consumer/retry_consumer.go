package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aakash811/queueflow/shared/tracing"
	"github.com/aakash811/queueflow/worker-service/models"
	kafkago "github.com/segmentio/kafka-go"
)

func StartRetryConsumer(jobChannel chan jobWithContext) {
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

		headers := make(map[string]string)
		for _, h := range message.Headers {
			headers[h.Key] = string(h.Value)
		}

		traceCtx := tracing.ExtractTraceContext(
			context.Background(),
			headers,
		)

		publishTime := time.Now()
		if pt, ok := headers["publish-time"]; ok {
			if t, err := time.Parse(time.RFC3339Nano, pt); err == nil {
				publishTime = t
			}
		}

		jobChannel <- jobWithContext{Job: job, Ctx: traceCtx, PublishTime: publishTime}
	}
}