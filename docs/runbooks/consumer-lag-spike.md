# Consumer Lag Spike

## Symptoms
- Grafana dashboard shows consumer lag increasing rapidly
- Jobs are not being picked up by workers
- `queue_depth` metric is growing

## Diagnosis
1. Check worker pod status: `kubectl get pods -n queueflow -l app=worker-service`
2. Check worker logs: `kubectl logs -l app=worker-service -n queueflow --tail=100`
3. Check Kafka topic partitions: `kafka-topics --describe --topic jobs_pending --bootstrap-server $KAFKA_BROKERS`
4. Verify consumer group assignment: `kafka-consumer-groups --describe --group queueflow-workers --bootstrap-server $KAFKA_BROKERS`

## Resolution
- If workers are down: scale up `kubectl scale deployment worker-service --replicas=5 -n queueflow`
- If partition hot-spot: check partition key distribution; consider adding more partitions
- If broker is slow: check broker logs and disk I/O
- If dead-letter queue is growing: check for poison messages

## Prevention
- Set up HPA with consumer lag metric (planned)
- Monitor broker health via Prometheus
- Configure max.poll.records appropriately
