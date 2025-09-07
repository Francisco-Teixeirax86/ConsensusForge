package metrics

import (
	"testing"
	"time"
)

func TestLabelHelpers(t *testing.T) {
	// Test NodeLabel
	nodeLabel := NodeLabel("node-1")
	if nodeLabel.Name != "node_id" || nodeLabel.Value != "node-1" {
		t.Errorf("NodeLabel failed: got %+v", nodeLabel)
	}
	
	// Test AlgorithmLabel
	algoLabel := AlgorithmLabel("raft")
	if algoLabel.Name != "algorithm" || algoLabel.Value != "raft" {
		t.Errorf("AlgorithmLabel failed: got %+v", algoLabel)
	}
	
	// Test StateLabel
	stateLabel := StateLabel("leader")
	if stateLabel.Name != "state" || stateLabel.Value != "leader" {
		t.Errorf("StateLabel failed: got %+v", stateLabel)
	}
	
	// Test CustomLabel
	customLabel := CustomLabel("custom", "value")
	if customLabel.Name != "custom" || customLabel.Value != "value" {
		t.Errorf("CustomLabel failed: got %+v", customLabel)
	}
}

func TestMetricConstants(t *testing.T) {
	// Test that metric constants are not empty
	constants := []string{
		MetricMessagesReceived,
		MetricMessagesSent,
		MetricElections,
		MetricLeaderChanges,
		MetricProposals,
		MetricCommittedEntries,
		MetricElectionDuration,
		MetricProposalDuration,
		MetricCommitLatency,
		MetricCurrentTerm,
		MetricLogSize,
		MetricCommitIndex,
		MetricActiveNodes,
	}
	
	for _, constant := range constants {
		if constant == "" {
			t.Error("Metric constant should not be empty")
		}
	}
	
	// Test for uniqueness
	seen := make(map[string]bool)
	for _, constant := range constants {
		if seen[constant] {
			t.Errorf("Duplicate metric constant: %s", constant)
		}
		seen[constant] = true
	}
}

func TestNoOpMetrics(t *testing.T) {
	metrics := NewNoOpMetrics()
	
	// These should not panic
	metrics.IncCounter("test")
	metrics.AddCounter("test", 1.0)
	metrics.SetGauge("test", 42.0)
	metrics.RecordHistogram("test", 1.5)
	metrics.RecordDuration("test", time.Second)
	
	// Test timer
	timer := metrics.StartTimer("test")
	if timer == nil {
		t.Error("StartTimer should return a timer")
	}
	timer.Stop() // Should not panic
}
