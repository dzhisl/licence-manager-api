package middleware

import (
	"strconv"

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

	ErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_errors_total",
			Help: "Total number of error requests processed by the server.",
		},
		[]string{"path", "status"},
	)
)

func PrometheusInit() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(ErrorCount)
}

func TrackMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
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

	}
}
