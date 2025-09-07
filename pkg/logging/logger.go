package logging

import (
	"context"
	"time"
)

// LogLevel represents the severity of a log message
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func (l LogLevel) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger defines the interface for structured logging
type Logger interface {
	// Debug logs a debug message with optional fields
	Debug(msg string, fields ...Field)

	// Info logs an info message with optional fields
	Info(msg string, fields ...Field)

	// Warn logs a warning message with optional fields
	Warn(msg string, fields ...Field)

	// Error logs an error message with optional fields
	Error(msg string, fields ...Field)

	// Fatal logs a fatal message and exits
	Fatal(msg string, fields ...Field)

	// With returns a new logger with additional context fields
	With(fields ...Field) Logger

	// WithContext returns a new logger with context
	WithContext(ctx context.Context) Logger
}

// Field represents a structured logging field
type Field struct {
	Key   string
	Value interface{}
}

// Helper functions for creating fields
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Value: value}
}
