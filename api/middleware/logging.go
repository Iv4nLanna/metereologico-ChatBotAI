package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"weather-api/metrics"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()

		latencyMs := time.Since(start).Milliseconds()
		metrics.Default.RecordRequest(c.Writer.Status(), latencyMs)

		slog.Info("request",
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", latencyMs,
			"ip", c.ClientIP(),
		)
	}
}
