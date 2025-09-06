package consensus

import (
	"context"

	"github.com/francisco-teixeirax86/consensusforge/pkg/config"
)

// Defines the interface for consensus algorithm implementations
type Algorithm interface {
	Name() string

	CreateNode(id string, config config.Config) (Node, error)
}

// Represents a consensus algorithm node
type Node interface {
	Start(ctx context.Context) error
	Stop() error
	ID() string
	IsLeader() bool
	Propose(data []byte) error
	GetState() NodeState
}

// Represents the current state of a consensus node
type NodeState int

const (
	StateFollower NodeState = iota
	StateCandidate
	StateLeader
	StateStopped
)

func (s NodeState) String() string {
	switch s {
	case StateFollower:
		return "Follower"
	case StateCandidate:
		return "Candidate"
	case StateLeader:
		return "Leader"
	case StateStopped:
		return "Stopped"
	default:
		return "Unknown"
	}
}
