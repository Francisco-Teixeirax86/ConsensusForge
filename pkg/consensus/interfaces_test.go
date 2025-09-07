package consensus

import (
	"testing"
)

func TestNodeStateString(t *testing.T) {
	tests := []struct {
		state    NodeState
		expected string
	}{
		{StateFollower, "Follower"},
		{StateCandidate, "Candidate"},
		{StateLeader, "Leader"},
		{StateStopped, "Stopped"},
		{NodeState(999), "Unknown"},
	}
	
	for _, test := range tests {
		result := test.state.String()
		if result != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, result)
		}
	}
}

func TestMessageType(t *testing.T) {
	// Test that message types are distinct
	types := []MessageType{
		MessageAppendEntries,
		MessageRequestVote,
		MessageAppendEntriesResponse,
		MessageRequestVoteResponse,
		MessagePrepare,
		MessagePromise,
		MessageAccept,
		MessageAccepted,
		MessageHeartbeat,
		MessageClientRequest,
	}
	
	seen := make(map[MessageType]bool)
	for _, msgType := range types {
		if seen[msgType] {
			t.Errorf("Duplicate message type: %v", msgType)
		}
		seen[msgType] = true
	}
}

func TestEntryType(t *testing.T) {
	// Test that entry types are distinct
	types := []EntryType{
		EntryCommand,
		EntryConfig,
		EntrySnapshot,
	}
	
	seen := make(map[EntryType]bool)
	for _, entryType := range types {
		if seen[entryType] {
			t.Errorf("Duplicate entry type: %v", entryType)
		}
		seen[entryType] = true
	}
}
