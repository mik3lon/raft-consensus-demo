package raftnode

import (
	"encoding/json"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft"
)

type FSM struct {
	mu    sync.Mutex
	state map[string]string
}

func NewFSM() *FSM {
	return &FSM{
		state: make(map[string]string),
	}
}

func (f *FSM) Apply(log *raft.Log) interface{} {
	f.mu.Lock()
	defer f.mu.Unlock()

	var command map[string]string
	if err := json.Unmarshal(log.Data, &command); err != nil {
		fmt.Printf("Failed to unmarshal command: %v", err)
		return "ERROR"
	}

	for key, value := range command {
		f.state[key] = value
	}
	return "SUCCESS"
}

func (f *FSM) Snapshot() (raft.FSMSnapshot, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	snapshot := make(map[string]string)
	for key, value := range f.state {
		snapshot[key] = value
	}
	return &FSMSnapshot{state: snapshot}, nil
}

func (f *FSM) Restore(data io.ReadCloser) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	defer data.Close()
	return json.NewDecoder(data).Decode(&f.state)
}

type FSMSnapshot struct {
	state map[string]string
}

func (s *FSMSnapshot) Persist(sink raft.SnapshotSink) error {
	err := func() error {
		data, err := json.Marshal(s.state)
		if err != nil {
			return err
		}
		if _, err := sink.Write(data); err != nil {
			return err
		}
		return sink.Close()
	}()
	if err != nil {
		sink.Cancel()
	}
	return err
}

func (s *FSMSnapshot) Release() {}
