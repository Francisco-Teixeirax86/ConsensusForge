package metrics

import "time"

// Metrics defines the interface for collecting consensus algorithm metrics
type Metrics interface {
	// Counter metrics
	IncCounter(name string, labels ...Label)
	AddCounter(name string, value float64, labels ...Label)

	// Gauge metrics
	SetGauge(name string, value float64, labels ...Label)

	// Histogram metrics
	RecordHistogram(name string, value float64, labels ...Label)

	// Timer convenience method
	RecordDuration(name string, duration time.Duration, labels ...Label)

	// Create a timer that records when stopped
	StartTimer(name string, labels ...Label) Timer
}

// Timer represents a running timer
type Timer interface {
	// Stop the timer and record the duration
	Stop()
}

// Label represents a metric label
type Label struct {
	Name  string
	Value string
}

// Helper functions for creating labels
func NodeLabel(nodeID string) Label {
	return Label{Name: "node_id", Value: nodeID}
}

func AlgorithmLabel(algorithm string) Label {
	return Label{Name: "algorithm", Value: algorithm}
}

func StateLabel(state string) Label {
	return Label{Name: "state", Value: state}
}

func CustomLabel(name, value string) Label {
	return Label{Name: name, Value: value}
}

// Common metric names as constants
const (
	// Consensus metrics
	MetricMessagesReceived = "consensus_messages_received_total"
	MetricMessagesSent     = "consensus_messages_sent_total"
	MetricElections        = "consensus_elections_total"
	MetricLeaderChanges    = "consensus_leader_changes_total"
	MetricProposals        = "consensus_proposals_total"
	MetricCommittedEntries = "consensus_committed_entries_total"

	// Performance metrics
	MetricElectionDuration = "consensus_election_duration_seconds"
	MetricProposalDuration = "consensus_proposal_duration_seconds"
	MetricCommitLatency    = "consensus_commit_latency_seconds"

	// State metrics
	MetricCurrentTerm = "consensus_current_term"
	MetricLogSize     = "consensus_log_size"
	MetricCommitIndex = "consensus_commit_index"
	MetricActiveNodes = "consensus_active_nodes"
)
