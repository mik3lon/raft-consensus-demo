package raftnode

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/raft"
)

func HandleSet(r *raft.Raft, body io.ReadCloser, w http.ResponseWriter) {
	defer body.Close()

	var command map[string]string
	if err := json.NewDecoder(body).Decode(&command); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	data, _ := json.Marshal(command)
	future := r.Apply(data, 10*time.Second)
	if err := future.Error(); err != nil {
		http.Error(w, "Failed to apply command: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("SUCCESS"))
}

func HandleGet(fsm *FSM, w http.ResponseWriter) {
	fsm.mu.Lock()
	defer fsm.mu.Unlock()

	data := fsm.state
	json.NewEncoder(w).Encode(data)
}

func HandleLeader(r *raft.Raft, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	leader := r.Leader()
	json.NewEncoder(w).Encode(map[string]string{
		"leader": string(leader),
	})
}
