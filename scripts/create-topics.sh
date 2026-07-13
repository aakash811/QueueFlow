#!/usr/bin/env bash
set -euo pipefail

PARTITIONS="${KAFKA_DEFAULT_PARTITIONS:-8}"

topics=(
  jobs_pending
  jobs_processing
  jobs_completed
  jobs_failed
  jobs_retry
  jobs_deadletter
)

for topic in "${topics[@]}"
do
  kafka-topics --create \
    --if-not-exists \
    --topic "$topic" \
    --bootstrap-server "${KAFKA_BROKERS:-localhost:9092}" \
    --partitions "$PARTITIONS" \
    --replication-factor 1
done