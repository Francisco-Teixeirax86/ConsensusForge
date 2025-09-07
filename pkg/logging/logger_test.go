package logging

import (
	"context"
	"testing"
	"time"
)

func TestLogLevelString(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{LevelFatal, "FATAL"},
		{LogLevel(999), "UNKNOWN"},
	}
	
	for _, test := range tests {
		result := test.level.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestFieldHelpers(t *testing.T) {
	// Test String field
	field := String("key", "value")
	if field.Key != "key" || field.Value != "value" {
		t.Errorf("String field failed: got %+v", field)
	}
	
	// Test Int field
	intField := Int("count", 42)
	if intField.Key != "count" || intField.Value != 42 {
		t.Errorf("Int field failed: got %+v", intField)
	}
	
	// Test Int64 field
	int64Field := Int64("bignum", 123456789)
	if int64Field.Key != "bignum" || int64Field.Value != int64(123456789) {
		t.Errorf("Int64 field failed: got %+v", int64Field)
	}
	
	// Test Duration field
	duration := 5 * time.Second
	durationField := Duration("timeout", duration)
	if durationField.Key != "timeout" || durationField.Value != duration {
		t.Errorf("Duration field failed: got %+v", durationField)
	}
	
	// Test Any field
	anyField := Any("data", map[string]int{"test": 1})
	if anyField.Key != "data" {
		t.Errorf("Any field failed: got %+v", anyField)
	}
}

func TestNoOpLogger(t *testing.T) {
	logger := NewNoOpLogger()
	
	// These should not panic
	logger.Debug("test")
	logger.Info("test")
	logger.Warn("test")
	logger.Error("test")
	
	// Test chaining
	newLogger := logger.With(String("key", "value"))
	if newLogger == nil {
		t.Error("With() should return a logger")
	}
	
	ctxLogger := logger.WithContext(context.Background())
	if ctxLogger == nil {
		t.Error("WithContext() should return a logger")
	}
}
