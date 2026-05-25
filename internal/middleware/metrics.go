package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MetricsMiddleware struct {
	requests map[string]int64
	latency  map[string][]time.Duration
}

func NewMetricsMiddleware() *MetricsMiddleware {
	return &MetricsMiddleware{
		requests: make(map[string]int64),
		latency:  make(map[string][]time.Duration),
	}
}

func (m *MetricsMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		path := c.FullPath()
		if path == "" {
			path = "unknown"
		}
		method := c.Request.Method
		status := strconv.Itoa(c.Writer.Status())
		key := method + ":" + path + ":" + status

		m.requests[key]++

		if m.latency[key] == nil {
			m.latency[key] = make([]time.Duration, 0, 1000)
		}
		m.latency[key] = append(m.latency[key], duration)
	}
}

func (m *MetricsMiddleware) GetRequests() map[string]int64 {
	return m.requests
}

func (m *MetricsMiddleware) GetAvgLatency(path string) time.Duration {
	durations := m.latency[path]
	if len(durations) == 0 {
		return 0
	}
	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}