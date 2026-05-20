package models

import "time"

type DeadLetterJob struct {
	ID          string                 `json:"id"`
	QueueName   string                 `json:"queue_name"`
	Payload     map[string]interface{} `json:"payload"`
	RetryCount  int                    `json:"retry_count"`
	FailedAt    time.Time              `json:"failed_at"`
}