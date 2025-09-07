package network

import (
	"testing"
)

func TestNetworkManagerCreation(t *testing.T) {
	manager := NewNetworkManager()
	
	if manager.transports == nil {
		t.Error("Expected transports map to be initialized")
	}
	
	nodes := manager.GetAllNodes()
	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes initially, got %d", len(nodes))
	}
}

func TestNetworkManagerCreateRemoveNode(t *testing.T) {
	manager := NewNetworkManager()
	
	// Create node
	transport := manager.CreateNode("node-1")
	if transport == nil {
		t.Fatal("CreateNode should return a transport")
	}
	
	nodes := manager.GetAllNodes()
	if len(nodes) != 1 {
		t.Errorf("Expected 1 node, got %d", len(nodes))
	}
	
	if nodes[0] != "node-1" {
		t.Errorf("Expected node 'node-1', got '%s'", nodes[0])
	}
	
	// Get node
	retrieved, err := manager.GetNode("node-1")
	if err != nil {
		t.Fatalf("GetNode failed: %v", err)
	}
	if retrieved != transport {
		t.Error("Retrieved transport should be the same as created")
	}
	
	// Remove node
	err = manager.RemoveNode("node-1")
	if err != nil {
		t.Fatalf("RemoveNode failed: %v", err)
	}
	
	nodes = manager.GetAllNodes()
	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes after removal, got %d", len(nodes))
	}
	
	// Try to get removed node
	_, err = manager.GetNode("node-1")
	if err == nil {
		t.Error("GetNode should fail for removed node")
	}
}

func TestNetworkManagerMultipleNodes(t *testing.T) {
	manager := NewNetworkManager()
	
	// Create multiple nodes
	nodeIDs := []string{"node-1", "node-2", "node-3"}
	transports := make(map[string]NetworkTransport)
	
	for _, nodeID := range nodeIDs {
		transport := manager.CreateNode(nodeID)
		transports[nodeID] = transport
	}
	
	// Check all nodes exist
	nodes := manager.GetAllNodes()
	if len(nodes) != 3 {
		t.Errorf("Expected 3 nodes, got %d", len(nodes))
	}
	
	// Verify each node can be retrieved
	for _, nodeID := range nodeIDs {
		transport, err := manager.GetNode(nodeID)
		if err != nil {
			t.Errorf("GetNode failed for %s: %v", nodeID, err)
		}
		if transport != transports[nodeID] {
			t.Errorf("Retrieved transport doesn't match for %s", nodeID)
		}
	}
}

func TestNetworkManagerShutdown(t *testing.T) {
	manager := NewNetworkManager()
	
	// Create some nodes
	manager.CreateNode("node-1")
	manager.CreateNode("node-2")
	
	nodes := manager.GetAllNodes()
	if len(nodes) != 2 {
		t.Errorf("Expected 2 nodes, got %d", len(nodes))
	}
	
	// Shutdown
	err := manager.Shutdown()
	if err != nil {
		t.Fatalf("Shutdown failed: %v", err)
	}
	
	// Check all nodes are gone
	nodes = manager.GetAllNodes()
	if len(nodes) != 0 {
		t.Errorf("Expected 0 nodes after shutdown, got %d", len(nodes))
	}
}
