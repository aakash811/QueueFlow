# QueueFlow — Distributed Job Processing Platform

## Overview

QueueFlow is a distributed asynchronous job processing platform built using Go, Kafka, PostgreSQL, Redis, gRPC, Prometheus, and Grafana.

The system simulates production-grade backend infrastructure capable of handling asynchronous workflows, distributed worker orchestration, retries, dead-letter queues, delayed jobs, recurring cron jobs, observability, authentication, rate limiting, and scalable event-driven processing.

QueueFlow demonstrates modern backend engineering concepts used in real-world distributed systems and infrastructure platforms.

---

# Features

- Distributed asynchronous job processing
- Kafka-based event-driven architecture
- Concurrent worker pool orchestration
- Retry mechanisms with exponential backoff
- Dead Letter Queue (DLQ)
- Delayed job scheduling
- Recurring cron jobs
- JWT authentication
- Role-Based Access Control (RBAC)
- Redis-based rate limiting
- Backpressure handling
- Circuit breaker pattern
- Prometheus metrics
- Grafana dashboards
- gRPC internal communication
- NGINX API Gateway
- CI/CD using GitHub Actions
- Load testing using k6
- Kafka partition optimization
- Worker concurrency tuning

---

# Architecture Diagram

![QueueFlow Architecture](assets/architecture-diagram.png)

---

# System Architecture

QueueFlow follows a distributed event-driven microservice architecture.

### Core Components

| Component         | Responsibility                          |
| ----------------- | --------------------------------------- |
| Job Service       | Handles API requests and publishes jobs |
| Kafka             | Distributed event queue                 |
| Worker Service    | Consumes and processes jobs             |
| Scheduler Service | Handles delayed and recurring jobs      |
| PostgreSQL        | Persistent job storage                  |
| Redis             | Rate limiting and idempotency           |
| Prometheus        | Metrics collection                      |
| Grafana           | Monitoring dashboards                   |
| NGINX             | Reverse proxy and API gateway           |
| gRPC              | Internal service communication          |

---

# Tech Stack

## Backend

- Go (Golang)
- Gin Framework
- gRPC
- Zap Logger

## Infrastructure

- Apache Kafka
- PostgreSQL
- Redis
- Docker
- Docker Compose
- NGINX

## Observability

- Prometheus
- Grafana
- OpenTelemetry distributed tracing

## Tracing

QueueFlow v2 instruments every job lifecycle span end-to-end using OpenTelemetry:

- **Job Service**: spans cover HTTP request → Kafka publish → DB write
- **Worker Service**: spans cover Kafka consume → job processing → DB status update → retry/DLQ publish
- **Scheduler Service**: spans cover DB poll → Kafka publish for delayed/cron jobs

Trace context is propagated across service boundaries via Kafka message headers using W3C TraceContext.

### Viewing Traces

Jaeger UI is included in `docker compose`:

1. Start the stack: `docker compose up --build`
2. Open Jaeger: http://localhost:16686
3. Select a service (e.g. `job-service`) and click **Search traces**

Traces show the full lifecycle of a single job across API, Kafka, and worker.

## Security

- JWT Authentication
- RBAC
- Rate Limiting

## Performance Testing

### k6 Load Testing

QueueFlow includes k6 load tests for measuring throughput scaling with Kafka partition count.

#### Quick Single Test
```bash
k6 run loadtest/job_creation.js
```

#### Partition Scaling Benchmark

Test how throughput scales with partition count:

```bash
# Git Bash / WSL
export PARTITIONS_LIST="1,2,4,8,16"
bash loadtest/run.sh

# Windows PowerShell
.\loadtest\run_benchmark.ps1 -Partitions "1,2,4,8,16"
```

#### Metrics Extracted

After each test run, these metrics are captured from Prometheus:

| Metric | Description |
|--------|-------------|
| `jobs_created_total` | Total jobs submitted |
| `job_throughput_total` | Total jobs completed |
| `queue_depth` | Current queue depth |
| `kafka_consumer_lag` | Unprocessed messages |
| `job_latency_seconds` | End-to-end processing latency |
| `job_pickup_latency_seconds` | Time from publish to worker pickup |
| `job_completion_rate` | Ratio of completed to created jobs |

#### Sample Resume Metrics

After running the partition scaling benchmark, you can report:

- **Throughput:** ~X req/s at 8 partitions (near-linear scaling from 1→8 partitions)
- **P95 Latency:** <X ms end-to-end
- **P99 Latency:** <X ms end-to-end
- **Job Completion Rate:** 99.9%+ (zero job loss)
- **Consumer Lag:** Bounded at <X messages under load

#### Grafana Dashboard

View real-time metrics at http://localhost:3000

Dashboard panels:
- Job Throughput
- Queue Depth
- Job Retry
- Job Failures
- Job Latency (P95/P99)
- Job Pickup Latency (P95/P99)

---

# Distributed Workflow

## Job Processing Flow

```text
Client Request
↓
NGINX API Gateway
↓
Job Service
↓
Kafka Topic
↓
Worker Service
↓
PostgreSQL Status Update
```

---

# Retry & Dead Letter Queue (DLQ)

QueueFlow implements fault-tolerant processing using retries and dead-letter queues.

## Retry Mechanism

- Exponential backoff retry strategy
- Configurable retry count
- Retry metrics tracking

## Dead Letter Queue

Failed jobs exceeding retry limits are:

- published to DLQ
- persisted in PostgreSQL
- logged for observability

---

# Scheduling System

QueueFlow supports:

## Delayed Jobs

```json
{
  "execute_at": "2026-05-20T10:00:00Z"
}
```

## Recurring Cron Jobs

Implemented using:

- robfig/cron

Scheduler service dispatches recurring jobs to Kafka asynchronously.

---

# Security Features

## JWT Authentication

Secures:

- job APIs
- admin endpoints

## Role-Based Access Control

Supported roles:

- ADMIN
- USER
- WORKER

## Rate Limiting

Implemented using:

- Redis sliding window strategy

---

# Observability & Monitoring

QueueFlow provides production-grade observability using Prometheus and Grafana.

## Metrics Tracked

- Queue depth
- Job throughput
- Retry count
- Worker failures
- Job latency

---

# Monitoring Dashboards

## Queue Throughput

![Job Throughput](assets/job-throughput-graph.png)

---

## Queue Depth

![Queue Depth](assets/queue-depth-graph.png)

---

## Retry Metrics

![Retry Metrics](assets/job-retry-graph.png)

---

## Job Latency

![Job Latency](assets/job-latency-graph.png)

---

# gRPC Internal Communication

QueueFlow uses gRPC for:

- internal service communication
- health checks
- worker coordination
- status updates

---

# Performance Benchmarking

Load testing performed using k6.

## Methodology

Partition count is the primary lever for horizontal Kafka throughput. QueueFlow v2 routes jobs by `queue_name` partition key so that increasing partitions directly increases parallel consumer capacity.

Run the benchmark against a locally running stack while varying `KAFKA_DEFAULT_PARTITIONS`:

```bash
# Example: test with 1, 2, 4, 8, and 16 partitions
KAFKA_DEFAULT_PARTITIONS=1 docker compose up --build -d
k6 run -e KAFKA_PARTITIONS=1 loadtest/job_creation.js

KAFKA_DEFAULT_PARTITIONS=2 docker compose up --build -d
k6 run -e KAFKA_PARTITIONS=2 loadtest/job_creation.js

# repeat for 4, 8, 16...
```

Measure end-to-end throughput from job creation through worker completion using the worker-throughput metric exposed at `http://localhost:2112/metrics`.

## Benchmark Results

| Partitions | Workers | Avg Latency | Throughput  |
| ---------- | ------- | ----------- | ----------- |
| 1         | 4       | —           | —           |
| 2         | 4       | —           | —           |
| 4         | 4       | —           | —           |
| 8         | 4       | —           | —           |
| 16        | 4       | —           | —           |

_Replace dashes with measured values after running k6 benchmarks. Expected behavior: throughput scales near-linearly with partition count up to the tested ceiling (target: 500–2000 req/s)._

---

# Performance Optimizations

Implemented optimizations:

- Kafka partition tuning
- Worker concurrency tuning
- Database indexing
- Backpressure protection
- Circuit breaker pattern

---

# API Endpoints

## Authentication

```http
POST /login
```

## Create Job

```http
POST /jobs
```

## Delayed Job

```json
{
  "queue_name": "email",
  "payload": {
    "id": 1
  },
  "execute_at": "2026-05-20T10:00:00Z"
}
```

---

# Running Locally

## Clone Repository

```bash
git clone https://github.com/aakash811/QueueFlow.git
cd QueueFlow
```

## Start Services

```bash
docker compose up --build
```

## Database Migrations

Run migrations automatically via docker compose:

```bash
docker compose up migrate
```

Or run manually:

```bash
docker compose up -d postgres
migrate -path scripts/migrations -database "postgres://postgres:postgres@localhost:5432/queueflow?sslmode=disable" up
```

## Authentication

### Default Admin User
- Username: `admin`
- Password: `admin123` (bcrypt hashed in database)

### API Keys
Generate an API key for tenant-based rate limiting:

```bash
go run scripts/generate_api_key.go <tenant_id> <key_name>
```

Include the API key in requests:
```bash
curl -H "X-API-Key: <your-api-key>" http://localhost:8080/jobs
```

## Secrets Management

Local development uses `.env` files. Production deployments use AWS Secrets Manager:

1. Create a secret in AWS Secrets Manager with JSON payload:
```json
{
  "POSTGRES_URL": "postgres://...",
  "REDIS_URL": "redis://...",
  "KAFKA_BROKERS": "kafka:9092",
  "JWT_SECRET": "...",
  "DB_PASSWORD": "..."
}
```

2. Set environment variable:
```bash
export AWS_SECRETS_MANAGER_SECRET_ID="arn:aws:secretsmanager:..."
```

3. Restart services - they will fetch secrets from AWS Secrets Manager instead of `.env`.

---

# Services

| Service     | Port  |
| ----------- | ----- |
| Job Service | 8080  |
| Prometheus  | 9090  |
| Grafana     | 3000  |
| Jaeger      | 16686 |
| Metrics     | 2112  |
| gRPC        | 50051 |

---

# Future Improvements

- Persistent cron orchestration
- Chaos testing automation

---

# Cloud Deployment

QueueFlow v2 includes Terraform modules and Kubernetes manifests for deploying to AWS EKS.

## Prerequisites

- AWS CLI configured with credentials
- kubectl installed
- Terraform >= 1.0

## Deploy Infrastructure

```bash
cd infra
terraform init
terraform plan
terraform apply
```

This provisions:
- VPC with public/private subnets
- EKS cluster with managed node group
- RDS PostgreSQL 16
- MSK Kafka cluster

## Deploy Application

```bash
# Configure kubectl
aws eks update-kubeconfig --name queueflow-cluster --region us-east-1

# Deploy services
./deployments/deploy.sh
```

## Verify Deployment

```bash
kubectl get pods -n queueflow
kubectl logs -l app=job-service -n queueflow
kubectl logs -l app=worker-service -n queueflow
```

---

# Key Engineering Concepts Demonstrated

- Distributed systems
- Event-driven architecture
- Asynchronous processing
- Fault tolerance
- Scalability
- Observability
- Infrastructure engineering
- Microservices
- Performance optimization
- Backend reliability engineering

---

# Author

Aakash Borse

GitHub:
https://github.com/aakash811
