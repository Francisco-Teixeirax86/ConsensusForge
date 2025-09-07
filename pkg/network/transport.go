package network

import (
	"time"

	"github.com/francisco-teixeirax86/consensusforge/pkg/consensus"
)

type NetworkConditions struct {
	BaseLatency   time.Duration `yaml:"base_latency"`
	LatencyJitter time.Duration `yaml:"latency_jitter"`

	PacketLoss  float64 `yaml:"packet_loss"` // 0.0 to 1.0
	Duplication float64 `yaml:"duplication"` // 0.0 to 1.0

	Bandwidth int64 `yaml:"bandwidth"` // bytes per second

	Partitioned    bool     `yaml:"partitioned"`
	PartitionNodes []string `yaml:"partition_nodes"`

	Corruption float64 `yaml:"corruption"` // 0.0 to 1.0
}

// DefaultNetworkConditions returns normal network conditions
func DefaultNetworkConditions() NetworkConditions {
	return NetworkConditions{
		BaseLatency:    1 * time.Millisecond,
		LatencyJitter:  0,
		PacketLoss:     0.0,
		Duplication:    0.0,
		Bandwidth:      0, // unlimited
		Partitioned:    false,
		PartitionNodes: []string{},
		Corruption:     0.0,
	}
}

// Extends the basic transport with network simulation capabilities
type NetworkTransport interface {
	consensus.Transport

	SetConditions(from, to string, conditions NetworkConditions)
	GetConditions(from, to string) NetworkConditions

	CreatePartition(nodes []string) error
	RemovePartition(nodes []string) error
	ClearPartitions() error

	GetStats() NetworkStats
	ResetStats()
}

// Contains network statistics
type NetworkStats struct {
	MessagesSent       int64                `json:"messages_sent"`
	MessagesReceived   int64                `json:"messages_received"`
	MessagesDropped    int64                `json:"messages_dropped"`
	MessagesDuplicated int64                `json:"messages_duplicated"`
	MessagesCorrupted  int64                `json:"messages_corrupted"`
	AverageLatency     time.Duration        `json:"average_latency"`
	NodeStats          map[string]NodeStats `json:"node_stats"`
}

// Contains per-node statistics
type NodeStats struct {
	Sent     int64 `json:"sent"`
	Received int64 `json:"received"`
	Dropped  int64 `json:"dropped"`
}
