package ports

import "context"

// LoggerBag is a struct that holds key-value pairs for logging.
type LoggerBag struct {
	Val interface{}
	Key string
}

// NewLoggerBag creates a new LoggerBag instance.
func NewLoggerBag(key string, val interface{}) LoggerBag {
	return LoggerBag{
		Key: key,
		Val: val,
	}
}

// LoggingService is an interface that defines the methods for logging operations.
type LoggingService interface {
	Info(ctx context.Context, msg string, bag ...LoggerBag)
	Debug(ctx context.Context, msg string, bag ...LoggerBag)
	Warn(ctx context.Context, msg string, bag ...LoggerBag)
	Error(ctx context.Context, msg string, bag ...LoggerBag)
}
