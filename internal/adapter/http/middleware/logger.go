package middleware

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lapotkin/file-storage/internal/logging"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = io.Discard
}

// Logger logs HTTP requests using the custom logger, skipping the /health endpoint
func Logger() gin.HandlerFunc {
	return gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	})
}

// LoggerWithCustomLogger logs HTTP requests using the provided custom logger
func LoggerWithCustomLogger(logger *logging.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		entry := logger.WithField("status", statusCode).
			WithField("latency", latency.String()).
			WithField("client_ip", clientIP).
			WithField("method", method).
			WithField("path", path)

		switch {
		case statusCode >= 500:
			entry.Error("HTTP Request")
		case statusCode >= 400:
			entry.Warn("HTTP Request")
		default:
			entry.Info("HTTP Request")
		}
	}
}
