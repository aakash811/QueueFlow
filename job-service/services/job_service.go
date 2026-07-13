package services

import (
	"context"
	"fmt"

	"github.com/aakash811/queueflow/job-service/kafka"
	"github.com/aakash811/queueflow/job-service/models"
	"github.com/aakash811/queueflow/job-service/repository"
	"github.com/aakash811/queueflow/shared/logger"
	"github.com/aakash811/queueflow/shared/metrics"
	"github.com/aakash811/queueflow/shared/tracing"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func CreateJob(ctx context.Context, job models.Job) error {
	job.ID = uuid.New().String()
	job.Status = "pending"
	job.RetryCount = 0
	job.MaxRetries = 3

	err := repository.CreateJob(job)

	if err != nil {
		return err
	}

	metrics.JobsCreated.Inc()

	if job.ExecuteAt != nil {
		fmt.Println(
			"Scheduled job stored for later execution:",
			job.ID,
		)

		return nil
	}

	logger.Log.Info( 
		"publishing immediate job", 
		
		zap.String("job_id", job.ID), 
		zap.String("queue", job.QueueName), 
	)
	
	ctx, span := tracing.Tracer.Start(ctx, "publish-job")
	defer span.End()

	err = kafka.PublishJob(ctx, job)
	
	if err != nil {
		return err
	}
	metrics.QueueDepth.Inc()
	
	return nil
}