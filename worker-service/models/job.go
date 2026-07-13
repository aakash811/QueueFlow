package models

import "time"

type Job struct {
	ID          string                 `json:"id"`
	QueueName   string                 `json:"queue_name"`
	Payload     map[string]interface{} `json:"payload"`
	Status      string                 `json:"status"`
	RetryCount  int                    `json:"retry_count"`
	MaxRetries  int                    `json:"max_retries"`
	CreatedAt   *time.Time             `json:"created_at"`
	UpdatedAt   *time.Time             `json:"updated_at"`
	ProcessedAt *time.Time             `json:"processed_at"`
	FailedAt    *time.Time             `json:"failed_at"`
	ExecuteAt   *time.Time             `json:"execute_at"`
	TenantID    string                 `json:"tenant_id"`
	PartitionKey string                `json:"partition_key"`
	Priority    int                    `json:"priority"`
}