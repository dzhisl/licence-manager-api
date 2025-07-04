package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of requests processed by the server.",
		},
		[]string{"path", "status"},
	)
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "status"},
	)

	ErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_errors_total",
			Help: "Total number of error requests processed by the server.",
		},
		[]string{"path", "status"},
	)

	SuccessCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_success_total",
			Help: "Total number of request with 200 status",
		},
		[]string{"path", "status"},
	)
)

func PrometheusInit() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(ErrorCount)
	prometheus.MustRegister(SuccessCount)
	prometheus.MustRegister(RequestDuration)
}

func TrackMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		status := strconv.Itoa(c.Writer.Status())

		RequestCount.WithLabelValues(path, status).Inc()

		// Optionally log errors if status >= 400
		if c.Writer.Status() >= 400 {
			ErrorCount.WithLabelValues(path, status).Inc()
		}
		if c.Writer.Status() == 200 || c.Writer.Status() == 201 {
			SuccessCount.WithLabelValues(path, status).Inc()
		}
		duration := time.Since(start).Seconds()
		RequestDuration.WithLabelValues(path, status).Observe(duration)

	}
}
