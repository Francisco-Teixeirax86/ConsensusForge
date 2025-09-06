package consensus

type StateMachine interface {
	Apply(data []byte) ([]byte, error)

	Snapshot() ([]byte, error)

	Restore(snapshot []byte) error

	GetState() interface{}
}

// Represents a log entry in the consensus log
type Entry struct {
	Index   int64     `json:"index"`
	Term    int64     `json:"term"`
	Command []byte    `json:"command"`
	Type    EntryType `json:"type"`
}

// Defines de type of log entry
type EntryType int

const (
	EntryCommand EntryType = iota
	EntryConfig
	EntrySnapshot
)
