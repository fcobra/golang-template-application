package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace = "base_app"
)

// Metrics holds the Prometheus metrics collectors.
type Metrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	RequestsInFlight prometheus.Gauge
}

// New creates and registers the metrics.
func New(appName string) *Metrics {
	m := &Metrics{
		RequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "Total number of HTTP requests.",
				ConstLabels: prometheus.Labels{
					"app": appName,
				},
			},
			[]string{"method", "path", "status_code"},
		),
		RequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "Duration of HTTP requests.",
				ConstLabels: prometheus.Labels{
					"app": appName,
				},
				Buckets: prometheus.DefBuckets, // Default buckets
			},
			[]string{"method", "path"},
		),
		RequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "requests_in_flight",
				Help:      "Number of current in-flight HTTP requests.",
				ConstLabels: prometheus.Labels{
					"app": appName,
				},
			},
		),
	}
	return m
}
