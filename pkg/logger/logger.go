package logger

import (
	"context"

	"github.com/dzhisl/license-api/internal/api/middleware"
	"github.com/dzhisl/license-api/pkg/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// skipLevel is the number of stack frames to ascend to report the correct caller
const skipLevel = 1

var log zap.Logger

// InitLogger sets up the global logger based on the environment
func InitLogger() {
	var cfg zap.Config
	var logLevel string
	environment := config.AppConfig.StageLevel

	// Use JSON logger for production, console logger for development
	if environment == "production" {
		cfg = zap.NewProductionConfig()
		cfg.EncoderConfig.TimeKey = "timestamp"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		logLevel = "debug"
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// Disable stacktrace to reduce verbosity
	cfg.EncoderConfig.StacktraceKey = ""

	// Set log level from configuration

	switch logLevel {
	case "debug":
		cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	default:
		cfg.Level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	}

	logger, err := cfg.Build()
	if err != nil {
		// If we can't build the logger, use a default logger to report the error
		zap.NewExample().Fatal("Error building logger", zap.Error(err))
	}

	// Replace the global logger
	log = *logger
	zap.ReplaceGlobals(logger)
}

func loggerMiddleware(ctx context.Context, fields []zap.Field) []zap.Field {
	fields = extractReqId(ctx, fields)
	// in future there may be more middlewares here
	return fields
}

func extractReqId(ctx context.Context, fields []zap.Field) []zap.Field {
	reqID, ok := ctx.Value(middleware.RequestIDKey).(string)
	if !ok {
		return fields
	}
	if reqID != "" {
		newField := zap.String(string(middleware.RequestIDKey), reqID)
		fields = append(fields, newField)
	}
	return fields
}

// Info logs an info message with optional fields
func Info(ctx context.Context, msg string, fields ...zap.Field) {
	fields = loggerMiddleware(ctx, fields)
	log.WithOptions(zap.AddCallerSkip(skipLevel)).Info(msg, fields...)
}

// Error logs an info message with optional fields
func Error(ctx context.Context, msg string, fields ...zap.Field) {
	fields = loggerMiddleware(ctx, fields)
	log.WithOptions(zap.AddCallerSkip(skipLevel)).Error(msg, fields...)
}

// Debug logs a debug message with optional fields
func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	fields = loggerMiddleware(ctx, fields)
	log.WithOptions(zap.AddCallerSkip(skipLevel)).Debug(msg, fields...)
}

// Warn logs a warning message with optional fields
func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	fields = loggerMiddleware(ctx, fields)
	log.WithOptions(zap.AddCallerSkip(skipLevel)).Warn(msg, fields...)
}

// Fatal logs a fatal message with optional fields and then exits
func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	fields = loggerMiddleware(ctx, fields)
	log.WithOptions(zap.AddCallerSkip(skipLevel)).Fatal(msg, fields...)
}
