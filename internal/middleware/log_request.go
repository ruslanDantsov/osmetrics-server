package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func NewLoggerRequestMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		responseSize := c.Writer.Size()

		logger.Info("HTTP Request",
			zap.String("uri", c.Request.URL.Path),
			zap.String("method", c.Request.Method),
			zap.String("duration", duration.String()),
			zap.Int("status", status),
			zap.Int("response_size_bytes", responseSize),
		)
	}
}
