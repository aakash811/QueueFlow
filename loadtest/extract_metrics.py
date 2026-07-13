#!/usr/bin/env python3
"""Extract Prometheus metrics for QueueFlow benchmark reporting."""

import json
import sys
import urllib.request

PROM_URL = sys.argv[1] if len(sys.argv) > 1 else "http://localhost:9090"

queries = {
    "jobs_created_total": 'jobs_created_total',
    "job_throughput_total": 'job_throughput_total',
    "queue_depth": 'queue_depth',
    "kafka_consumer_lag": 'kafka_consumer_lag',
    "job_latency_p95": 'histogram_quantile(0.95, rate(job_latency_seconds_bucket[5m]))',
    "job_latency_p99": 'histogram_quantile(0.99, rate(job_latency_seconds_bucket[5m]))',
    "job_pickup_p95": 'histogram_quantile(0.95, rate(job_pickup_latency_seconds_bucket[5m]))',
    "job_pickup_p99": 'histogram_quantile(0.99, rate(job_pickup_latency_seconds_bucket[5m]))',
    "job_completion_rate": 'sum(rate(job_throughput_total[5m])) / sum(rate(jobs_created_total[5m]))',
    "worker_failures_total": 'worker_failures_total',
}

results = {}
for name, query in queries.items():
    url = f"{PROM_URL}/api/v1/query?query={urllib.parse.quote(query)}"
    try:
        with urllib.request.urlopen(url, timeout=10) as resp:
            data = json.loads(resp.read().decode())
            if data["data"]["result"]:
                results[name] = data["data"]["result"][0]["value"][1]
            else:
                results[name] = "0"
    except Exception as e:
        results[name] = f"error: {e}"

for k, v in results.items():
    print(f"{k}: {v}")
