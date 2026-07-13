# Hot Partition

## Symptoms
- One Kafka partition has significantly higher lag than others
- Some workers are idle while others are overloaded
- Latency spikes correlate with specific partition IDs

## Diagnosis
1. Check partition distribution: `kafka-topics --describe --topic jobs_pending --bootstrap-server $KAFKA_BROKERS`
2. Identify partition key causing hotspot: `SELECT partition_key, COUNT(*) FROM jobs GROUP BY partition_key ORDER BY COUNT(*) DESC;`
3. Check if a single tenant/queue is dominating traffic

## Resolution
- Short term: increase partition count for affected topic
- Medium term: review partition key strategy; add salt/jitter to key
- Long term: implement per-tenant rate limiting

## Prevention
- Use composite partition keys (e.g., `tenant_id:queue_name`)
- Monitor per-partition lag in Grafana
- Set up alerts when partition lag variance exceeds threshold
