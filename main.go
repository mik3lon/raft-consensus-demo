package main

import (
	"github.com/mik3lon/raft-consensus-demo/raftnode"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("Usage: go run main.go <node_id> <raft_port> <http_port>")
	}

	nodeID := os.Args[1]
	raftPort := os.Args[2]
	httpPort := os.Args[3]

	// Initialize Raft node
	node, err := raftnode.NewRaftNode(nodeID, raftPort)
	if err != nil {
		log.Fatalf("Failed to create raft node: %v", err)
	}

	// Start HTTP server for client interactions
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		raftnode.HandleSet(node.Raft, r.Body, w)
	})
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		raftnode.HandleGet(node.FSM, w)
	})
	http.HandleFunc("/leader", func(w http.ResponseWriter, r *http.Request) {
		raftnode.HandleLeader(node.Raft, w)
	})

	log.Printf("Node %s listening for client requests on port %s", nodeID, httpPort)
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}
