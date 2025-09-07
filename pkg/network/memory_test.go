package network

import (
	"testing"
	"time"

	"github.com/francisco-teixeirax86/consensusforge/pkg/consensus"
)

func TestDefaultNetworkConditions(t *testing.T) {
	conditions := DefaultNetworkConditions()
	
	if conditions.BaseLatency != 1*time.Millisecond {
		t.Errorf("Expected BaseLatency 1ms, got %v", conditions.BaseLatency)
	}
	
	if conditions.PacketLoss != 0.0 {
		t.Errorf("Expected PacketLoss 0.0, got %f", conditions.PacketLoss)
	}
	
	if conditions.Partitioned {
		t.Error("Expected Partitioned false")
	}
}

func TestMemoryTransportCreation(t *testing.T) {
	transport := NewMemoryTransport("node-1")
	
	if transport.nodeID != "node-1" {
		t.Errorf("Expected nodeID 'node-1', got '%s'", transport.nodeID)
	}
	
	if transport.inbox == nil {
		t.Error("Expected inbox to be initialized")
	}
	
	if transport.nodes == nil {
		t.Error("Expected nodes map to be initialized")
	}
	
	if transport.conditions == nil {
		t.Error("Expected conditions map to be initialized")
	}
}

func TestMemoryTransportSendReceive(t *testing.T) {
	// Create two nodes
	node1 := NewMemoryTransport("node-1")
	node2 := NewMemoryTransport("node-2")
	
	// Connect them
	nodes := map[string]*MemoryTransport{
		"node-1": node1,
		"node-2": node2,
	}
	node1.Connect(nodes)
	node2.Connect(nodes)
	
	// Send message from node1 to node2
	msg := consensus.Message{
		Type:      consensus.MessageHeartbeat,
		From:      "node-1",
		To:        "node-2",
		Data:      []byte("hello"),
		Timestamp: time.Now(),
	}
	
	err := node1.Send("node-2", msg)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	
	// Receive message at node2
	select {
	case received := <-node2.Receive():
		if received.From != msg.From {
			t.Errorf("Expected From '%s', got '%s'", msg.From, received.From)
		}
		if string(received.Data) != string(msg.Data) {
			t.Errorf("Expected Data '%s', got '%s'", string(msg.Data), string(received.Data))
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Message not received within timeout")
	}
}

func TestMemoryTransportBroadcast(t *testing.T) {
	// Create three nodes
	node1 := NewMemoryTransport("node-1")
	node2 := NewMemoryTransport("node-2")
	node3 := NewMemoryTransport("node-3")
	
	// Connect them
	nodes := map[string]*MemoryTransport{
		"node-1": node1,
		"node-2": node2,
		"node-3": node3,
	}
	node1.Connect(nodes)
	node2.Connect(nodes)
	node3.Connect(nodes)
	
	// Broadcast from node1
	msg := consensus.Message{
		Type:      consensus.MessageHeartbeat,
		From:      "node-1",
		Data:      []byte("broadcast"),
		Timestamp: time.Now(),
	}
	
	err := node1.Broadcast(msg)
	if err != nil {
		t.Fatalf("Broadcast failed: %v", err)
	}
	
	// Check that both node2 and node3 received the message
	received := 0
	timeout := time.After(100 * time.Millisecond)
	
	for received < 2 {
		select {
		case <-node2.Receive():
			received++
		case <-node3.Receive():
			received++
		case <-timeout:
			t.Fatalf("Only received %d/2 messages", received)
		}
	}
}

func TestMemoryTransportPacketLoss(t *testing.T) {
	node1 := NewMemoryTransport("node-1")
	node2 := NewMemoryTransport("node-2")
	
	nodes := map[string]*MemoryTransport{
		"node-1": node1,
		"node-2": node2,
	}
	node1.Connect(nodes)
	node2.Connect(nodes)
	
	// Set 100% packet loss
	conditions := DefaultNetworkConditions()
	conditions.PacketLoss = 1.0
	node1.SetConditions("node-1", "node-2", conditions)
	
	// Send message
	msg := consensus.Message{
		Type: consensus.MessageHeartbeat,
		From: "node-1",
		To:   "node-2",
		Data: []byte("lost"),
	}
	
	err := node1.Send("node-2", msg)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	
	// Should not receive message due to packet loss
	select {
	case <-node2.Receive():
		t.Fatal("Should not have received message due to packet loss")
	case <-time.After(50 * time.Millisecond):
		// Expected - message was lost
	}
	
	// Check stats
	stats := node1.GetStats()
	if stats.MessagesDropped == 0 {
		t.Error("Expected at least one dropped message")
	}
}

func TestMemoryTransportPartition(t *testing.T) {
	node1 := NewMemoryTransport("node-1")
	node2 := NewMemoryTransport("node-2")
	
	nodes := map[string]*MemoryTransport{
		"node-1": node1,
		"node-2": node2,
	}
	node1.Connect(nodes)
	node2.Connect(nodes)
	
	// Create partition
	err := node1.CreatePartition([]string{"node-1"})
	if err != nil {
		t.Fatalf("CreatePartition failed: %v", err)
	}
	
	// Send message - should be dropped due to partition
	msg := consensus.Message{
		Type: consensus.MessageHeartbeat,
		From: "node-1",
		To:   "node-2",
		Data: []byte("partitioned"),
	}
	
	err = node1.Send("node-2", msg)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	
	// Should not receive message due to partition
	select {
	case <-node2.Receive():
		t.Fatal("Should not have received message due to partition")
	case <-time.After(50 * time.Millisecond):
		// Expected - message was dropped
	}
	
	// Remove partition
	err = node1.RemovePartition([]string{"node-1"})
	if err != nil {
		t.Fatalf("RemovePartition failed: %v", err)
	}
	
	// Now message should go through
	err = node1.Send("node-2", msg)
	if err != nil {
		t.Fatalf("Send failed: %v", err)
	}
	
	select {
	case <-node2.Receive():
		// Expected - partition removed
	case <-time.After(50 * time.Millisecond):
		t.Fatal("Message should have been received after partition removal")
	}
}
