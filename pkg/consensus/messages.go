package consensus

import "time"

// Defines the type of consensus message
type MessageType int

const (
	// Raft message types
	MessageAppendEntries MessageType = iota
	MessageRequestVote
	MessageAppendEntriesResponse
	MessageRequestVoteResponse

	// Paxos message types
	MessagePrepare
	MessagePromise
	MessageAccept
	MessageAccepted

	// Generic types
	MessageHeartbeat
	MessageClientRequest
)

// Represents a consensus protocol message
type Message struct {
	Type      MessageType `json:"type"`
	From      string      `json:"from"`
	To        string      `json:"to"`
	Term      int64       `json:"term,omitempty"`
	Data      []byte      `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Defines the interface for message passing
type Transport interface {
	Send(to string, msg Message) error

	Broadcast(msg Message) error

	Receive() <-chan Message

	Close() error
}
