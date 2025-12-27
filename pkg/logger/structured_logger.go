package logger

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// StructuredLogger is the logging interface using zap
type StructuredLogger interface {
	Debug(ctx context.Context, message string, fields ...zap.Field)
	Info(ctx context.Context, message string, fields ...zap.Field)
	Warn(ctx context.Context, message string, fields ...zap.Field)
	Error(ctx context.Context, message string, err error, fields ...zap.Field)
	Fatal(ctx context.Context, message string, err error, fields ...zap.Field)
	WithRequestID(requestID string) StructuredLogger
	Sync() error
}

// zapLogger wraps zap.Logger to implement StructuredLogger interface
type zapLogger struct {
	logger    *zap.Logger
	requestID string
}

// NewStructuredLogger creates a new structured logger using zap
func NewStructuredLogger(logFilePath string, enableFile bool) (StructuredLogger, error) {
	var config zap.Config

	if enableFile && logFilePath != "" {
		// Production config with file output
		config = zap.NewProductionConfig()
		config.OutputPaths = []string{"stdout", logFilePath}
		config.ErrorOutputPaths = []string{"stderr", logFilePath}
	} else {
		// Development config with console output
		config = zap.NewDevelopmentConfig()
		config.OutputPaths = []string{"stdout"}
		config.ErrorOutputPaths = []string{"stderr"}
	}

	// Set format to JSON for better parsing
	config.Encoding = "json"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.TimeKey = "timestamp"

	logger, err := config.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build logger: %w", err)
	}

	return &zapLogger{logger: logger}, nil
}

// WithRequestID returns a new logger with request ID context
func (l *zapLogger) WithRequestID(requestID string) StructuredLogger {
	return &zapLogger{
		logger:    l.logger,
		requestID: requestID,
	}
}

// Debug logs a debug level message
func (l *zapLogger) Debug(ctx context.Context, message string, fields ...zap.Field) {
	fields = append(fields, zap.String("request_id", l.requestID))
	l.logger.Debug(message, fields...)
}

// Info logs an info level message
func (l *zapLogger) Info(ctx context.Context, message string, fields ...zap.Field) {
	fields = append(fields, zap.String("request_id", l.requestID))
	l.logger.Info(message, fields...)
}

// Warn logs a warning level message
func (l *zapLogger) Warn(ctx context.Context, message string, fields ...zap.Field) {
	fields = append(fields, zap.String("request_id", l.requestID))
	l.logger.Warn(message, fields...)
}

// Error logs an error level message with error details
func (l *zapLogger) Error(ctx context.Context, message string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	fields = append(fields, zap.String("request_id", l.requestID))
	l.logger.Error(message, fields...)
}

// Fatal logs a fatal level message and exits
func (l *zapLogger) Fatal(ctx context.Context, message string, err error, fields ...zap.Field) {
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	fields = append(fields, zap.String("request_id", l.requestID))
	l.logger.Fatal(message, fields...)
}

// Sync flushes any buffered log entries
func (l *zapLogger) Sync() error {
	return l.logger.Sync()
}

// Global logger instance
var globalLogger StructuredLogger

// InitializeGlobalLogger initializes the global logger instance
func InitializeGlobalLogger(logFilePath string, enableFile bool) error {
	logger, err := NewStructuredLogger(logFilePath, enableFile)
	if err != nil {
		return err
	}
	globalLogger = logger
	return nil
}

// GetLogger returns the global logger instance
func GetLogger() StructuredLogger {
	if globalLogger == nil {
		// Create a default logger if not initialized
		logger, _ := NewStructuredLogger("logs/app.log", true)
		globalLogger = logger
	}
	return globalLogger
}

// Convenience functions that use the global logger
func Debug(ctx context.Context, message string, fields ...zap.Field) {
	GetLogger().Debug(ctx, message, fields...)
}

func Info(ctx context.Context, message string, fields ...zap.Field) {
	GetLogger().Info(ctx, message, fields...)
}

func Warn(ctx context.Context, message string, fields ...zap.Field) {
	GetLogger().Warn(ctx, message, fields...)
}

func Error(ctx context.Context, message string, err error, fields ...zap.Field) {
	GetLogger().Error(ctx, message, err, fields...)
}

func Fatal(ctx context.Context, message string, err error, fields ...zap.Field) {
	GetLogger().Fatal(ctx, message, err, fields...)
}
