# QueueFlow v2 — Rebuild PRD

## 1. What's actually wrong with v1

Your own benchmark table is the tell: throughput stays flat at ~39 req/s whether you run 1 worker or 8. That means workers were never the bottleneck — the system has one serialization point somewhere upstream (almost certainly Kafka partition count, or a single-writer path into Postgres). A "distributed" system that doesn't get faster when you add distribution isn't distributed yet — it's a single pipeline with extra processes around it. That's the real gap v2 needs to close, and it's a better and more honest headline than "handles 39 req/s."

Everything below is designed around fixing that root cause, not just adding more feature checkboxes.

---

## 2. Problem Statement

Teams running background job processing (emails, webhooks, image processing, batch imports) need a system that: absorbs traffic spikes without dropping work, survives worker/broker failures without losing jobs, and scales horizontally in a way that's _provably_ linear — not just architecturally described as scalable.

v1 has the pieces (Kafka, workers, retries, DLQ) but never validated that scaling actually scales. v2's job is to prove it, and to close the operational gaps that separate a portfolio project from something you could hand to an on-call engineer.

## 3. Goals / Non-Goals

**Goals**

- Demonstrate horizontal scalability with data to back it up (throughput that actually increases with partitions/workers)
- Close the observability gap: distributed tracing, not just metrics
- Close the operations gap: infrastructure as code, one-command environment spin-up, defined SLOs
- Close the security gap: secrets management, dependency scanning, not just JWT+RBAC at the app layer

**Non-Goals**

- Not building a general-purpose message broker — Kafka does that job, you're building the orchestration layer around it
- Not targeting massive scale (millions of req/s) — targeting _provable, linear_ scale at a realistic ceiling (e.g. 500–2000 req/s) is more credible and more finishable than chasing a huge number you can't fully explain
- Not building a UI/dashboard product — Grafana + a minimal admin API is enough; don't burn time on frontend polish for a backend-signaling project

---

## 4. Functional Requirements

| #    | Requirement                                                                                                                                       |
| ---- | ------------------------------------------------------------------------------------------------------------------------------------------------- |
| FR1  | Accept job submissions via REST API with schema-validated payloads (job type, priority, payload, execute_at)                                      |
| FR2  | Route jobs to Kafka topics partitioned by a documented partition key strategy (e.g. tenant ID or job type) so parallelism is actually exploitable |
| FR3  | Support immediate, delayed, and recurring (cron) job execution                                                                                    |
| FR4  | Guarantee at-least-once delivery with idempotency keys so retries don't double-process                                                            |
| FR5  | Retry failed jobs with exponential backoff; route to DLQ after max attempts                                                                       |
| FR6  | Expose job status query API (pending / running / succeeded / failed / dead-lettered)                                                              |
| FR7  | Enforce JWT auth + RBAC (admin/user/worker) on all control-plane endpoints                                                                        |
| FR8  | Rate-limit per tenant/API key, not just globally                                                                                                  |
| FR9  | Provide an admin operation to requeue or purge DLQ items                                                                                          |
| FR10 | Support graceful worker shutdown (drain in-flight jobs before termination)                                                                        |

## 5. Non-Functional Requirements

**Scalability**

- Throughput must scale near-linearly with partition count up to a stated ceiling — this needs to be _measured and shown_, with a k6 report, not asserted
- Horizontal worker scaling via Kubernetes HPA based on consumer lag, not just CPU

**Reliability**

- Zero job loss under: broker restart, worker crash mid-processing, Postgres failover
- Defined RTO/RPO for the job store (even a simple stated target: "under 60s data loss window on Postgres failover via streaming replication")

**Observability**

- Metrics (have this — Prometheus/Grafana)
- **Distributed tracing** (missing today) — OpenTelemetry spans across API → Kafka → worker → DB, so you can trace one job's full lifecycle, not just aggregate counts
- Structured, correlation-ID-tagged logs across every service

**Security**

- Secrets via a vault/secret manager, not `.env` files in the repo pattern
- Dependency vulnerability scanning in CI (e.g. `govulncheck`, Trivy for images)
- Input validation and rate limiting at the gateway, not trusted to app code alone

**Operability**

- Infrastructure as code (Terraform or Pulumi) — currently `docker compose up` only, which doesn't demonstrate you can provision real infra
- One documented deployment path to a real cloud environment (even a small managed Kubernetes cluster), not just local Docker
- Runbooks: what to do when consumer lag spikes, when DLQ grows, when a partition is hot

---

## 6. What "production-ready" actually means here (checklist)

Use this as your definition of done — each unchecked item is a known gap, not a mystery:

- [ ] Kafka topics use a partition key that lets you demonstrate throughput scaling with partition count (the core fix)
- [ ] k6 benchmark report showing throughput vs. partition count as a curve, not a flat line
- [ ] OpenTelemetry distributed tracing wired end-to-end
- [ ] Idempotency keys implemented and tested (duplicate delivery doesn't double-process)
- [ ] Terraform/Pulumi provisions the infra (even to a single small cloud VM cluster)
- [ ] CI pipeline: lint, test, vulnerability scan, build, deploy — with branch protection
- [ ] Chaos test: kill a worker mid-job, kill the Kafka broker, kill Postgres primary — document what happens and confirm zero job loss
- [ ] SLOs defined and shown to be met (e.g. "P95 job pickup latency < 200ms", "99.9% jobs complete within 3 retries")
- [ ] Runbook doc for at least 3 failure scenarios
- [ ] Secrets not committed or `.env`-only — pulled from a secret manager at runtime

---

## 7. Suggested build order (so this stays finishable)

1. **Fix the partition/throughput story first.** This is the single highest-value fix — it directly answers the question your current README raises but doesn't resolve. Get a real scaling curve before anything else.
2. **Add tracing.** Second highest value — "I added OpenTelemetry tracing across the pipeline" is a strong, specific, verifiable claim.
3. **Add IaC + real cloud deploy.** This is what separates "I ran Docker Compose on my laptop" from "I can operate infrastructure."
4. **Chaos testing + runbooks.** This is the smallest time investment for a surprisingly large credibility gain — almost nobody at your level does this, and it's the kind of thing senior engineers specifically probe for.
5. Everything else in your existing "Future Improvements" list (multi-tenant queues, priority queues) is real but lower priority than the four above — they add features, not credibility.

---

## 8. The interview story this rebuild buys you

Right now the honest answer to "does this scale?" is "I don't know, my numbers didn't move." After this rebuild, the honest answer is: "I found the bottleneck was single-partition serialization, fixed the partition key strategy, and here's the curve showing throughput scaling with partition count up to N." That's a materially better technical story — it shows you can diagnose a system that _looks_ correct but isn't performing as designed, which is closer to real senior-engineer work than shipping more features on top of an unverified foundation.
