package metrics

import "time"

type MetricsFactory interface {
	CreateMetrics() Metrics

	Close() error
}

// Metrics implementation that does nothing
type NoOpMetrics struct {
}

func (n *NoOpMetrics) IncCounter(name string, labels ...Label) {
}

func (n *NoOpMetrics) AddCounter(name string, value float64, labels ...Label) {
}

func (n *NoOpMetrics) SetGauge(name string, value float64, labels ...Label) {
}

func (n *NoOpMetrics) RecordHistogram(name string, value float64, labels ...Label) {
}

func (n *NoOpMetrics) RecordDuration(name string, duration time.Duration, labels ...Label) {
}

func (n *NoOpMetrics) StartTimer(name string, labels ...Label) Timer {
	return &NoOpTimer{}
}

// Timer that does nothing
type NoOpTimer struct{}

func (n *NoOpTimer) Stop() {}

func NewNoOpMetrics() Metrics {
	return &NoOpMetrics{}
}
