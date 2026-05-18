```bash id="jlwm105"
#!/bin/bash

topics=(
  jobs.pending
  jobs.processing
  jobs.completed
  jobs.failed
  jobs.retry
  jobs.deadletter
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
```
