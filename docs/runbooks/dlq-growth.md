# DLQ Growth

## Symptoms
- Dead letter job count increasing in Grafana
- DLQ topic `jobs_deadletter` has growing lag
- Repeated failures for specific job types

## Diagnosis
1. Check DLQ size: `kafka-run-class kafka.tools.GetOffsetShell --topic jobs_deadletter --broker-list $KAFKA_BROKERS`
2. Query dead letter jobs: `SELECT * FROM dead_letter_jobs ORDER BY failed_at DESC LIMIT 20;`
3. Check failure reasons in logs
4. Identify if failures are systemic (DB down, external API timeout) or data-specific

## Resolution
- For systemic issues: fix the root cause, then requeue DLQ via admin API
- For data-specific issues: fix payload, then requeue
- Requeue command: `POST /admin/dlq/requeue` with job IDs
- Purge command: `POST /admin/dlq/purge` (use with caution)

## Prevention
- Set up alerts on DLQ growth rate
- Implement circuit breakers for external dependencies
- Add job payload schema validation
