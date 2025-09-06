package config

import (
	"time"
)

type Config struct {
	// Node configuration
	NodeID  string   `yaml:"node_id"`
	Peers   []string `yaml:"peers"`
	DataDir string   `yaml:"data_dir"`

	// Timing configuration
	ElectionTimeout   time.Duration `yaml:"election_timeout"`
	HeartbeatInterval time.Duration `yaml:"heartbeat_interval"`

	// Network configuration
	ListenAddr string `yaml:"listen_addr"`

	// Algorithm-specific settings
	Algorithm string                 `yaml:"algorithm"`
	Settings  map[string]interface{} `yaml:"settings,omitempty"`
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() Config {
	return Config{
		NodeID:            "node-1",
		Peers:             []string{},
		DataDir:           "./data",
		ElectionTimeout:   150 * time.Millisecond,
		HeartbeatInterval: 50 * time.Millisecond,
		ListenAddr:        ":8080",
		Algorithm:         "raft",
		Settings:          make(map[string]interface{}),
	}
}
