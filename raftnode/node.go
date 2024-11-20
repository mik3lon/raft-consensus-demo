package raftnode

import (
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/raft"
	"github.com/hashicorp/raft-boltdb"
)

type RaftNode struct {
	Raft *raft.Raft
	FSM  *FSM
}

// NewRaftNode initializes a Raft node with its transport and storage.
func NewRaftNode(nodeID, raftPort string) (*RaftNode, error) {
	// Raft configuration
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(nodeID)

	// Raft transport for internal communication
	addr := "127.0.0.1:" + raftPort
	transport, err := raft.NewTCPTransport(addr, nil, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to create transport: %w", err)
	}

	// Raft storage
	logStore, err := raftboltdb.NewBoltStore(nodeID + ".raft.db")
	if err != nil {
		return nil, fmt.Errorf("failed to create log store: %w", err)
	}
	stableStore, err := raftboltdb.NewBoltStore(nodeID + ".stable.db")
	if err != nil {
		return nil, fmt.Errorf("failed to create stable store: %w", err)
	}
	snapshotStore, err := raft.NewFileSnapshotStore(".", 2, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot store: %w", err)
	}

	// FSM setup
	fsm := NewFSM()

	// Create Raft instance
	r, err := raft.NewRaft(config, fsm, logStore, stableStore, snapshotStore, transport)
	if err != nil {
		return nil, fmt.Errorf("failed to create raft: %w", err)
	}

	// Bootstrap cluster only for node1
	if nodeID == "node1" {
		config := raft.Configuration{
			Servers: []raft.Server{
				{ID: raft.ServerID("node1"), Address: raft.ServerAddress("127.0.0.1:8081")},
				{ID: raft.ServerID("node2"), Address: raft.ServerAddress("127.0.0.1:8082")},
				{ID: raft.ServerID("node3"), Address: raft.ServerAddress("127.0.0.1:8083")},
			},
		}
		if err := r.BootstrapCluster(config).Error(); err != nil && err != raft.ErrCantBootstrap {
			return nil, fmt.Errorf("failed to bootstrap cluster: %w", err)
		}
	}

	return &RaftNode{
		Raft: r,
		FSM:  fsm,
	}, nil
}
