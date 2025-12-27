package api

import "github.com/prometheus/client_golang/prometheus"

var requestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_ms",
		Help:    "HTTP request duration in milliseconds",
		Buckets: []float64{5, 10, 20, 35, 50, 75, 100, 150, 200, 400, 800, 1600},
	},
	[]string{"method", "route", "status"},
)

func RegisterMetrics() {
	prometheus.MustRegister(requestDuration)
}
