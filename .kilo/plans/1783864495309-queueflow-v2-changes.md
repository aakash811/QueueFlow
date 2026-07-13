# QueueFlow v2 Rebuild Plan

## Current State
- 3 Go services: job-service (Gin API), worker-service (Kafka consumer pool), scheduler-service (delayed/cron)
- Kafka topics: 1 partition each → throughput flat at ~39 req/s regardless of worker count
- Auth: hardcoded admin/password, JWT secret in `.env` committed to repo
- Observability: Prometheus metrics only, no traces
- Deploy: `docker compose up` only
- CI: basic build-only GitHub Actions
- No IaC, no secret manager, no chaos tests, no runbooks, no SLOs
- Load test: k6 script exists but only measures submission latency, not full pipeline throughput

## PH1: Fix Partition Key & Throughput Scaling (Highest Priority)

**Decision: CHOSEN** — Partition key = `queue_name`.  
Rationale: simple, requires no new tenant domain model, gives good distribution across queues, aligns with existing queue abstraction. If multi-tenancy is added later, key can become `tenant_id:queue_name`.

Tasks:
1. Update `scripts/create-topics.sh` — create all topics with configurable partition count (default 8, env-driven)
2. Update `docker-compose.yml` — add `KAFKA_NUM_PARTITIONS` env or use `kafka-configs` to set defaults
3. Update all Kafka producers (`job-service/kafka/producer.go`, `worker-service/kafka/producer.go`, `scheduler-service/kafka/producer.go`) — set `Message.Key` to `queue_name`
4. Update worker consumer (`worker-service/consumer/consumer.go`) — already uses consumer group; verify partition assignment scales
5. Add config for `KafkaDefaultPartitions` in `shared/config/config.go`
6. Add k6 benchmark script that varies partition count (1, 2, 4, 8, 16) and plots throughput
7. Update README benchmark table with real scaling curve data and k6 report artifact
8. Fix bug in `worker-service/consumer/consumer.go:34` — `wg.Add(1)` is called inside the loop after `jobChannel <- job`, causing race; move to before channel send

## PH2: Distributed Tracing

**Decision: CHOSEN** — OpenTelemetry Go SDK with Jaeger all-in-one for dev, OTLP collector for prod.  
Rationale: standard, minimal code changes with auto-instrumentation, Jaeger UI is good enough for portfolio.

Tasks:
1. Add `go.opentelemetry.io/otel` packages to `go.mod`
2. Create `shared/tracing/tracing.go` — init tracer with resource attributes (service name, env)
3. Instrument job-service: Gin middleware auto-span + manual spans for Kafka publish + DB calls
4. Instrument worker-service: manual spans for Kafka consume + process + DB update + retry/DLQ publish
5. Instrument scheduler-service: spans for DB poll + Kafka publish
6. Propagate trace context via Kafka message headers using W3C traceparent
7. Add correlation ID to every structured log line (from span context)
8. Update `docker-compose.yml` — add Jaeger service on 16686
9. Add tracing section to README with Jaeger UI screenshot instructions

## PH3: Infrastructure as Code + Cloud Deploy

**Decision: CHOSEN** — Terraform.  
**Decision: NEEDS CLARIFICATION** — Target cloud: recommend AWS EKS or GCP GKE for managed K8s, but smallest credible target is acceptable. If no cloud preference, start with Terraform that provisions a single small VM cluster (e.g. DigitalOcean or AWS EC2) with K3s.

Tasks:
1. Add Terraform root module structure (`infra/`)
2. Provision VPC/network, K8s cluster (or small VM pool), Postgres (managed if possible), Kafka (managed or self-hosted on K8s)
3. Output kubeconfig and connection strings
4. Add `terraform.tfvars.example` and `.gitignore` for state
5. Add remote state backend config (S3 + DynamoDB or GCS)
6. Create Kubernetes manifests (`deployments/k8s/`) for all 3 services + HPA based on Kafka consumer lag
7. Add ArgoCD or Flux for GitOps (optional, good for portfolio)
8. Write one-command deploy script (`make deploy` or `./deploy.sh`)
9. Document cloud deployment path in README

**Decision: NEEDS CLARIFICATION** — Secrets manager: recommend HashiCorp Vault (self-hosted) or AWS Secrets Manager / GCP Secret Manager. Need user preference to pick exact implementation.

## PH4: Security & Operational Hardening

Tasks:
1. Remove `.env` from repo tracking; inject secrets via env at runtime
2. Integrate secret manager (Vault/AWS/GCP) into all 3 services — replace `config.go` viper env reads with secret fetcher (startup fetch or sidecar)
3. Add schema validation library (e.g. `github.com/go-playground/validator`) to job payloads
4. Replace hardcoded admin login with user table in Postgres + password hashing (bcrypt)
5. Change rate limiting from global to per-tenant API key (add `tenant_id` to Job model, generate API keys)
6. Add dependency vulnerability scanning to CI (`govulncheck` for Go, `trivy` for Docker images)
7. Add lint + unit tests to CI
8. Enable branch protection on main in GitHub settings (note: not Terraform-able via API easily, document manual step)
9. Verify graceful worker shutdown drains in-flight jobs before SIGTERM exit

## PH5: Observability SLOs + Chaos Testing + Runbooks

Tasks:
1. Define SLOs in `docs/slos.md`:
   - P95 job pickup latency < 200ms
   - P99 job pickup latency < 500ms  
   - 99.9% jobs complete within 3 retries
   - Consumer lag < 1000 messages under normal load
2. Add Grafana alerts for SLO burn rate
3. Write chaos tests (can be manual playbooks with `chaos-mesh` or just scripts):
   - Kill worker mid-job → confirm job is retried and completed
   - Kill Kafka broker → confirm zero message loss (acks=all)
   - Kill Postgres primary → confirm RTO/RPO < 60s (document streaming replication)
4. Write runbooks in `docs/runbooks/`:
   - `consumer-lag-spike.md`
   - `dlq-growth.md`
   - `hot-partition.md`

## PH6: Schema Migration + Data Model Fixes

Tasks:
1. Add `idempotency_key` column to `jobs` table + unique index
2. Add `tenant_id` column to `jobs` table + index
3. Add `partition_key` to `jobs` table for auditability
4. Add `priority` column (smallint) for future priority queue support
5. Add database migration tool (e.g. `golang-migrate`) instead of raw `init.sql`
6. Update all repository queries and models

## Validation Plan
1. Run `docker compose up` locally
2. Execute k6 benchmark with varying partition counts; confirm throughput scales near-linearly
3. Verify Jaeger UI shows complete traces for sample jobs
4. Run chaos tests and confirm zero job loss
5. Run CI pipeline and confirm lint + test + vuln scan + build all pass
6. Deploy to cloud staging via Terraform and run smoke tests

## Open Questions
1. **Cloud provider / cluster target**: AWS EKS, GCP GKE, or smallest self-hosted K8s (K3s on VMs)? Recommendation: GCP GKE smallest cluster or AWS EKS — lowest operational burden for portfolio.
2. **Secrets manager**: HashiCorp Vault vs AWS Secrets Manager vs GCP Secret Manager? Recommendation: AWS Secrets Manager if targeting AWS, GCP Secret Manager if targeting GCP — use cloud-native to reduce ops.
