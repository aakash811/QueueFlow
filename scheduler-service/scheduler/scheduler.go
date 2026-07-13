package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/aakash811/queueflow/scheduler-service/kafka"
	"github.com/aakash811/queueflow/scheduler-service/repository"
	"github.com/aakash811/queueflow/shared/tracing"
	"go.opentelemetry.io/otel/attribute"
)

func StartScheduler() {

	fmt.Println("scheduler service started")

	ticker := time.NewTicker(5 * time.Second)

	defer ticker.Stop()

	for {

		<-ticker.C

		ctx, span := tracing.Tracer.Start(context.Background(), "scheduler-poll")
		span.SetAttributes(attribute.String("correlation.id", tracing.GetCorrelationID(ctx)))

		jobs, err := repository.GetReadyJobs()

		if err != nil {

			fmt.Println(
				"scheduler fetch error:",
				err,
			)

			span.End()
			continue
		}

		for _, job := range jobs {

			jobCtx, jobSpan := tracing.Tracer.Start(ctx, "dispatch-job")
			jobSpan.SetAttributes(attribute.String("job.id", job.ID))

			fmt.Println(
				"dispatching scheduled job:",
				job.ID,
			)

			err = kafka.PublishJob(jobCtx, job)

			if err != nil {

				fmt.Println(
					"publish error:",
					err,
				)

				jobSpan.RecordError(err)
				jobSpan.End()
				continue
			}

			err = repository.MarkJobDispatched(
				job.ID,
			)

			if err != nil {

				fmt.Println(
					"mark dispatched error:",
					err,
				)
			}

			jobSpan.End()
		}

		span.End()
	}
}
