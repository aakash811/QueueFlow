package circuitbreaker

import (
	"fmt"
	"time"

	"github.com/sony/gobreaker"
)

var JobBreaker *gobreaker.CircuitBreaker

func InitCircuitBreaker() {
	settings := gobreaker.Settings{
		Name: "job-processor",
		MaxRequests: 3,
		Interval: 10 * time.Second,
		Timeout: 15 * time.Second,

		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 3
		},

		OnStateChange: func(name string, from, to gobreaker.State) {
			fmt.Printf(
				"circuit breaker state changed: %v -> %v\n",
				from,
				to,
			)
		},
	}

	JobBreaker = gobreaker.NewCircuitBreaker(
		settings,
	)
}