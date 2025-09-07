package network

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/francisco-teixeirax86/consensusforge/pkg/consensus"
)

// Implements the NetworkTransport for in-memory testing
type MemoryTransport struct {
	nodeID     string
	nodes      map[string]*MemoryTransport
	inbox      chan consensus.Message
	conditions map[string]NetworkConditions // from -> to
	partitions map[string]bool              // partitioned nodes
	stats      NetworkStats
	mu         sync.RWMutex
	closed     bool
}

// Creates a new in-memory transport
func NewMemoryTransport(nodeID string) *MemoryTransport {
	return &MemoryTransport{
		nodeID:     nodeID,
		nodes:      make(map[string]*MemoryTransport),
		inbox:      make(chan consensus.Message, 1000),
		conditions: make(map[string]NetworkConditions),
		partitions: make(map[string]bool),
		stats: NetworkStats{
			NodeStats: make(map[string]NodeStats),
		},
	}
}

// Connects this transport to other nodes
func (mt *MemoryTransport) Connect(nodes map[string]*MemoryTransport) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.nodes = nodes
	// Set default conditions for all nodes
	for nodeID := range nodes {
		if nodeID != mt.nodeID {
			key := fmt.Sprintf("%s->%s", mt.nodeID, nodeID)
			mt.conditions[key] = DefaultNetworkConditions()
		}
	}
}

func (mt *MemoryTransport) Send(to string, msg consensus.Message) error {
	mt.mu.RLock()
	if mt.closed {
		mt.mu.RUnlock()
		return fmt.Errorf("transport closed")
	}

	target, exists := mt.nodes[to]
	if !exists {
		mt.mu.RUnlock()
		return fmt.Errorf("node %s not found", to)
	}

	// Check if nodes are partitioned
	if mt.partitions[mt.nodeID] != mt.partitions[to] &&
		(mt.partitions[mt.nodeID] || mt.partitions[to]) {
		mt.stats.MessagesDropped++
		mt.mu.RUnlock()
		return nil // silently drop
	}

	key := fmt.Sprintf("%s->%s", mt.nodeID, to)
	conditions := mt.conditions[key]
	mt.mu.RUnlock()

	// Aplly network conditions
	if rand.Float64() < conditions.PacketLoss {
		mt.mu.Lock()
		mt.stats.MessagesDropped++
		mt.mu.Unlock()
		return nil // packet loss
	}

	// Calculate latency
	latency := conditions.BaseLatency
	if conditions.LatencyJitter > 0 {
		jitter := time.Duration(rand.Int64N(int64(conditions.LatencyJitter)))
		latency += jitter
	}

	go func() {
		if latency > 0 {
			time.Sleep(latency)
		}

		finalMsg := msg
		if rand.Float64() < conditions.Duplication {
			finalMsg.Data = []byte("corrupted")
			mt.mu.Lock()
			mt.stats.MessagesCorrupted++
			mt.mu.Unlock()
		}

		// Deliver message
		select {
		case target.inbox <- finalMsg:
			mt.mu.Lock()
			mt.stats.MessagesSent++
			nodeStats := mt.stats.NodeStats[to]
			nodeStats.Sent++
			mt.stats.NodeStats[to] = nodeStats
			mt.mu.Unlock()

		default:
			// inbox is full, drop message
			mt.mu.Lock()
			mt.stats.MessagesDropped++
			mt.mu.Unlock()
		}

		if rand.Float64() < conditions.Duplication {
			select {
			case target.inbox <- finalMsg:
				mt.mu.Lock()
				mt.stats.MessagesDuplicated++
				mt.mu.Unlock()
			default:
			}
		}
	}()

	return nil
}

// Broadcast Implements the consensus.Transport
func (mt *MemoryTransport) Broadcast(msg consensus.Message) error {
	mt.mu.RLock()
	if mt.closed {
		mt.mu.RUnlock()
		return fmt.Errorf("transport closed")
	}

	nodes := make([]string, 0, len(mt.nodes))
	for nodeID := range mt.nodes {
		if nodeID != mt.nodeID {
			nodes = append(nodes, nodeID)
		}
	}
	mt.mu.RUnlock()

	for _, nodeID := range nodes {
		if err := mt.Send(nodeID, msg); err != nil {
			return err
		}
	}
	return nil
}

// Reviece Implements the consensus.Transport
func (mt *MemoryTransport) Receive() <-chan consensus.Message {
	return mt.inbox
}

// Close Implements the consensus.Transport
func (mt *MemoryTransport) Close() error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	if !mt.closed {
		close(mt.inbox)
		mt.closed = true
	}
	return nil
}

func (mt *MemoryTransport) SetConditions(from, to string, conditions NetworkConditions) {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	key := fmt.Sprintf("%s->%s", from, to)
	mt.conditions[key] = conditions
}

func (mt *MemoryTransport) GetConditions(from, to string) NetworkConditions {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	key := fmt.Sprintf("%s->%s", from, to)
	if conditions, exists := mt.conditions[key]; exists {
		return conditions
	}
	return DefaultNetworkConditions()
}

func (mt *MemoryTransport) CreatePartition(nodes []string) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	for _, nodeID := range nodes {
		mt.partitions[nodeID] = true
	}
	return nil
}

func (mt *MemoryTransport) RemovePartition(nodes []string) error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	for _, nodeID := range nodes {
		delete(mt.partitions, nodeID)
	}
	return nil
}

func (mt *MemoryTransport) ClearPartitions() error {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.partitions = make(map[string]bool)
	return nil
}

func (mt *MemoryTransport) GetStats() NetworkStats {
	mt.mu.RLock()
	defer mt.mu.RUnlock()

	stats := mt.stats
	stats.NodeStats = make(map[string]NodeStats)
	for key, value := range mt.stats.NodeStats {
		stats.NodeStats[key] = value
	}
	return stats
}

func (mt *MemoryTransport) ResetStats() {
	mt.mu.Lock()
	defer mt.mu.Unlock()

	mt.stats = NetworkStats{
		NodeStats: make(map[string]NodeStats),
	}
}
