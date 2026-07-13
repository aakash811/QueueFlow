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

	JobsCreated = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "jobs_created_total",
			Help: "Total jobs created",
		},
	)

	JobLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "job_latency_seconds",
			Help: "Job processing latency",
			Buckets: prometheus.DefBuckets,
		},
	)

	JobPickupLatency = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "job_pickup_latency_seconds",
			Help: "Time from Kafka publish to worker start-of-processing",
			Buckets: prometheus.DefBuckets,
		},
	)

	KafkaConsumerLag = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "kafka_consumer_lag",
			Help: "Current Kafka consumer lag for jobs_pending topic",
		},
	)
)

func InitMetrics() {
	prometheus.MustRegister(
		QueueDepth,
		JobThroughput,
		WorkerFailures,
		RetryCount,
		JobsCreated,
		JobLatency,
		JobPickupLatency,
		KafkaConsumerLag,
	)
}