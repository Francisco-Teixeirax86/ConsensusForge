package logging

import "context"

// LoggerFactory creates logger instances
type LoggerFactory interface {
	// CreateLogger creates a new logger with the given name
	CreateLogger(name string) Logger

	// SetLevel sets the global log level
	SetLevel(level LogLevel)

	// Close shuts down the logging system
	Close() error
}

// NoOpLogger is a logger that does nothing (useful for testing)
type NoOpLogger struct{}

func (n *NoOpLogger) Debug(msg string, fields ...Field)      {}
func (n *NoOpLogger) Info(msg string, fields ...Field)       {}
func (n *NoOpLogger) Warn(msg string, fields ...Field)       {}
func (n *NoOpLogger) Error(msg string, fields ...Field)      {}
func (n *NoOpLogger) Fatal(msg string, fields ...Field)      {}
func (n *NoOpLogger) With(fields ...Field) Logger            { return n }
func (n *NoOpLogger) WithContext(ctx context.Context) Logger { return n }

// NewNoOpLogger returns a no-op logger
func NewNoOpLogger() Logger {
	return &NoOpLogger{}
}
