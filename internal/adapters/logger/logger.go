// Package logger provides a custom logger implementation using the slog package.
package logger

import (
	"context"
	"log/slog"
	"os"

	port "gogs.utking.net/utking/spaces/internal/ports"
)

var (
	logLevelMap = map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"INFO":  slog.LevelInfo,
		"WARN":  slog.LevelWarn,
		"ERROR": slog.LevelError,
	}
)

// AppLogger is a custom logger that wraps the slog.Logger and provides additional functionality.
type AppLogger struct {
	logger *slog.Logger
}

// NewAdapter creates a new instance of AppLogger with the specified log file.
func NewAdapter(logFile *os.File, level string) *AppLogger {
	logger := slog.New(slog.NewJSONHandler(
		logFile, &slog.HandlerOptions{
			Level: logLevelMap[level],
		}))

	return &AppLogger{
		logger: logger,
	}
}

// Info logs an informational message with the given key-value pairs.
func (l *AppLogger) Info(ctx context.Context, msg string, bag ...port.LoggerBag) {
	l.logger.InfoContext(ctx, msg, bagToAny(bag...)...)
}

// bagToAny converts a slice of LoggerBag to a slice of any.
func bagToAny(bag ...port.LoggerBag) []any {
	anyList := make([]any, 0, len(bag)*2)

	for _, b := range bag {
		anyList = append(anyList, b.Key, b.Val)
	}

	return anyList
}

// Debug logs a debug message with the given key-value pairs.
func (l *AppLogger) Debug(ctx context.Context, msg string, bag ...port.LoggerBag) {
	l.logger.DebugContext(ctx, msg, bagToAny(bag...)...)
}

// Warn logs a warning message with the given key-value pairs.
func (l *AppLogger) Warn(ctx context.Context, msg string, bag ...port.LoggerBag) {
	l.logger.WarnContext(ctx, msg, bagToAny(bag...)...)
}

// Error logs an error message with the given key-value pairs.
func (l *AppLogger) Error(ctx context.Context, msg string, bag ...port.LoggerBag) {
	l.logger.ErrorContext(ctx, msg, bagToAny(bag...)...)
}
