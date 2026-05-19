package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	QueueDepth = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "queue_depth",
			Help: "Current queue depth",
		},
	)

	JobThroughput = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "job_throughput_total",
			Help: "Total processed jobs",
		},
	)

	WorkerFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "worker_failures_total",
			Help: "Total worker failures",
		},
	)

	RetryCount = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "job_retry_total",
			Help: "Total retried jobs",
		},
	)

	JobLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "job_latency_seconds",
			Help: "Job processing latency",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func InitMetrics() {
	prometheus.MustRegister(
		QueueDepth,
		JobThroughput,
		WorkerFailures,
		RetryCount,
		JobLatency,
	)
}