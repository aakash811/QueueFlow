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
    --bootstrap-server localhost:9092 \
    --partitions 1 \
    --replication-factor 1
done