# Raft Consensus Demo

This project is a simple implementation of a Raft consensus algorithm using Go. It provides a minimal working example
for learning and experimenting with distributed consensus mechanisms.

## Features

- **Raft Consensus Implementation**: A leader election and log replication protocol for managing a replicated state
  machine.
- **HashiCorp Raft Library**: Uses the robust Raft implementation provided by HashiCorp.
- **HTTP Interface**: Simple HTTP endpoints to interact with the Raft cluster for setting and retrieving values.

---

## Table of Contents

1. [Getting Started](#getting-started)
2. [Usage](#usage)
3. [Architecture](#architecture)
4. [Endpoints](#endpoints)

---

## Getting Started

### Prerequisites

- Go 1.18 or higher
- A terminal to run multiple instances of the application

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/mik3lon/raft-consensus-demo.git
   cd raft-consensus-demo
    ```

2. Install dependencies:

    ```bash
    go mod tidy
    ```

3. Compile the application:

    ```bash
    go build -o raft-demo main.go
    ```

### Usage

Run the application with the following command:

    ```bash
    ./raft-demo <node_id> <raft_port> <http_port>
    ```

#### Example

Start three nodes:

    ```bash
    ./raft-demo node1 8081 9001
    ./raft-demo node2 8082 9002
    ./raft-demo node3 8083 9003
    Once the nodes are up, they form a Raft cluster with node1 as the initial leader (if bootstrap conditions are met).
    ```

### Architecture
#### Raft Node Initialization
Each node initializes:

* Raft Transport: Handles communication between nodes.
* Storage:
    * Log Store: Persistent storage for Raft logs.
    * Stable Store: Persistent storage for stable state (e.g., term and votedFor).
    * Snapshot Store: Persistent snapshots of the state machine.
* FSM (Finite State Machine): Application state replicated across all nodes.

### Cluster Bootstrapping
The cluster is bootstrapped by node1, which predefines the initial set of servers.

## Endpoints

| **Endpoint** | **Method** | **Description**                                   | **Example Request**                                        |
|--------------|------------|---------------------------------------------------|------------------------------------------------------------|
| `/set`       | POST       | Sets a value in the FSM.                          | `curl -X POST -d '{"key": "example", "value": "raft"}' http://localhost:9001/set` |
| `/get`       | GET        | Retrieves the current state from the FSM.         | `curl http://localhost:9001/get`                           |
| `/leader`    | GET        | Returns the current leader of the Raft cluster.   | `curl http://localhost:9001/leader`                        |
