package metrics

import "time"

// MetricsFactory creates metrics instances
type MetricsFactory interface {
	// CreateMetrics creates a new metrics instance
	CreateMetrics() Metrics

	// Close shuts down the metrics system
	Close() error
}

// NoOpMetrics is a metrics implementation that does nothing (useful for testing)
type NoOpMetrics struct{}

func (n *NoOpMetrics) IncCounter(name string, labels ...Label)                             {}
func (n *NoOpMetrics) AddCounter(name string, value float64, labels ...Label)              {}
func (n *NoOpMetrics) SetGauge(name string, value float64, labels ...Label)                {}
func (n *NoOpMetrics) RecordHistogram(name string, value float64, labels ...Label)         {}
func (n *NoOpMetrics) RecordDuration(name string, duration time.Duration, labels ...Label) {}
func (n *NoOpMetrics) StartTimer(name string, labels ...Label) Timer {
	return &NoOpTimer{}
}

// NoOpTimer is a timer that does nothing
type NoOpTimer struct{}

func (n *NoOpTimer) Stop() {}

// NewNoOpMetrics returns a no-op metrics implementation
func NewNoOpMetrics() Metrics {
	return &NoOpMetrics{}
}
