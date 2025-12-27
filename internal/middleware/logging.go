package middleware

import (
	"fmt"
	"ingredient-recognition-backend/pkg/logger"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	INFO = iota
	WARN
	ERROR
)

// LoggingMiddleware logs HTTP requests and responses
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request ID for tracing
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// Create a logger with the request ID
		ctx := c.Request.Context()
		loggerWithID := logger.GetLogger().WithRequestID(requestID)

		// Log request
		loggerWithID.Info(ctx, "Incoming request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("remote_addr", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// Record start time
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(startTime).Milliseconds()

		// Log response
		statusCode := c.Writer.Status()
		logLevel := INFO

		// Determine log level based on status code
		if statusCode >= 500 {
			logLevel = ERROR
		} else if statusCode >= 400 {
			logLevel = WARN
		} else if statusCode >= 200 && statusCode < 300 {
			logLevel = INFO
		}

		// Log based on level
		switch logLevel {
		case ERROR:
			loggerWithID.Error(ctx, "Request completed with error",
				nil,
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status_code", statusCode),
				zap.Int64("duration_ms", duration),
			)
		case WARN:
			loggerWithID.Warn(ctx, "Request completed with warning",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status_code", statusCode),
				zap.Int64("duration_ms", duration),
			)
		default:
			loggerWithID.Info(ctx, "Request completed successfully",
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.Int("status_code", statusCode),
				zap.Int64("duration_ms", duration),
			)
		}
	}
}

// ErrorHandlingMiddleware logs panics and recovers gracefully
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("request_id")
				loggerWithID := logger.GetLogger().WithRequestID(fmt.Sprint(requestID))
				loggerWithID.Error(c.Request.Context(), "Panic recovered",
					fmt.Errorf("%v", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
				)
				c.JSON(500, gin.H{"error": "Internal server error"})
			}
		}()
		c.Next()
	}
}
