package config

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	// Test default values
	if config.NodeID != "node-1" {
		t.Errorf("Expected NodeID 'node-1', got '%s'", config.NodeID)
	}
	
	if len(config.Peers) != 0 {
		t.Errorf("Expected empty peers slice, got %v", config.Peers)
	}
	
	if config.DataDir != "./data" {
		t.Errorf("Expected DataDir './data', got '%s'", config.DataDir)
	}
	
	if config.ElectionTimeout != 150*time.Millisecond {
		t.Errorf("Expected ElectionTimeout 150ms, got %v", config.ElectionTimeout)
	}
	
	if config.HeartbeatInterval != 50*time.Millisecond {
		t.Errorf("Expected HeartbeatInterval 50ms, got %v", config.HeartbeatInterval)
	}
	
	if config.ListenAddr != ":8080" {
		t.Errorf("Expected ListenAddr ':8080', got '%s'", config.ListenAddr)
	}
	
	if config.Algorithm != "raft" {
		t.Errorf("Expected Algorithm 'raft', got '%s'", config.Algorithm)
	}
	
	if config.Settings == nil {
		t.Error("Expected Settings to be initialized")
	}
}

func TestConfigModification(t *testing.T) {
	config := DefaultConfig()
	
	// Test modification
	config.NodeID = "test-node"
	config.Peers = []string{"peer1", "peer2"}
	config.ElectionTimeout = 200 * time.Millisecond
	
	if config.NodeID != "test-node" {
		t.Errorf("Expected NodeID 'test-node', got '%s'", config.NodeID)
	}
	
	if len(config.Peers) != 2 {
		t.Errorf("Expected 2 peers, got %d", len(config.Peers))
	}
	
	if config.ElectionTimeout != 200*time.Millisecond {
		t.Errorf("Expected ElectionTimeout 200ms, got %v", config.ElectionTimeout)
	}
}
