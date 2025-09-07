package network

import (
	"fmt"
	"sync"
)

type NetworkManager struct {
	transports map[string]*MemoryTransport
	mu         sync.RWMutex
}

func NewNetworkManager() *NetworkManager {
	return &NetworkManager{
		transports: make(map[string]*MemoryTransport),
	}
}

func (nm *NetworkManager) CreateNode(nodeID string) NetworkTransport {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	transport := NewMemoryTransport(nodeID)
	nm.transports[nodeID] = transport

	transport.Connect(nm.transports)

	for _, existingTrasport := range nm.transports {
		existingTrasport.Connect(nm.transports)
	}

	return transport
}

func (nm *NetworkManager) RemoveNode(nodeID string) error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	transport, exists := nm.transports[nodeID]
	if !exists {
		return fmt.Errorf("node %s not found", nodeID)
	}

	transport.Close()
	delete(nm.transports, nodeID)

	for _, existingTransport := range nm.transports {
		existingTransport.Connect(nm.transports)
	}

	return nil
}

func (nm *NetworkManager) GetNode(nodeID string) (NetworkTransport, error) {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	transport, exists := nm.transports[nodeID]
	if !exists {
		return nil, fmt.Errorf("node %s not found", nodeID)
	}

	return transport, nil
}

func (nm *NetworkManager) GetAllNodes() []string {
	nm.mu.RLock()
	defer nm.mu.RUnlock()

	nodes := make([]string, 0, len(nm.transports))
	for nodeID := range nm.transports {
		nodes = append(nodes, nodeID)
	}
	return nodes
}

func (nm *NetworkManager) Shutdown() error {
	nm.mu.Lock()
	defer nm.mu.Unlock()

	for _, transport := range nm.transports {
		transport.Close()
	}
	nm.transports = make(map[string]*MemoryTransport)
	return nil
}
