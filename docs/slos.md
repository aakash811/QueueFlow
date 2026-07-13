# QueueFlow SLOs

## Job Processing SLOs

| SLO | Target | Measurement Window |
|-----|--------|-------------------|
| P95 job pickup latency | < 200ms | 5 minutes |
| P99 job pickup latency | < 500ms | 5 minutes |
| Job completion rate (within 3 retries) | 99.9% | 24 hours |
| Consumer lag (normal load) | < 1,000 messages | 5 minutes |
| API availability | 99.9% | 30 days |

## SLO Definitions

### Job Pickup Latency
Time from job creation (Kafka publish) to worker start-of-processing.

Measured via: histogram `job_pickup_latency_seconds` in worker-service.

### Job Completion Rate
Percentage of jobs that reach `completed` status within `max_retries` (default: 3).

Measured via: `(completed_jobs / total_jobs) * 100` over 24h.

### Consumer Lag
Number of messages assigned to a consumer group but not yet processed.

Measured via: Kafka consumer lag metrics exposed to Prometheus.

## Alerting Thresholds

- **Warning**: P95 latency > 200ms for 5 minutes
- **Critical**: P95 latency > 500ms for 5 minutes
- **Warning**: Consumer lag > 1,000 for 5 minutes
- **Critical**: Consumer lag > 10,000 for 2 minutes
- **Warning**: Job completion rate < 99.5% over 1 hour
- **Critical**: Job completion rate < 99% over 1 hour

## Error Budget

- 24h error budget for 99.9% availability: 14.4s downtime per day, ~7.2m per month
- Burn rate alert triggers when budget is consumed faster than expected
