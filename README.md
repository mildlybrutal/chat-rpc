# RPC Chat Application

A real-time chat application built with Go's native RPC package, demonstrating distributed systems concepts with concurrent client-server architecture.


## Overview

This project implements a simple yet functional chat system using Go's built-in RPC (Remote Procedure Call) framework. Multiple clients can connect to a central server to send and receive messages in real-time.

## Architecture

```
┌─────────────┐    RPC Calls    ┌─────────────┐
│   Client    │ ◄─────────────► │   Server    │
│   ┌─────────┤                 ├─────────┐   │
│   │ Sender  │                 │ Message │   │
│   │ Thread  │                 │ Store   │   │
│   ├─────────┤                 ├─────────┤   │
│   │Receiver │                 │ Mutex   │   │
│   │ Thread  │                 │ Lock    │   │
│   └─────────┘                 └─────────┘   │
└─────────────┘                 └─────────────┘
```


### Components:
- **API Layer**: Shared interface definitions and data structures
- **Server**: RPC service handling message storage and retrieval
- **Client**: Concurrent sender/receiver with user interface

## How It Works

1. **Server** starts and listens on port 8080
2. **Clients** connect via TCP and register with username  
3. **Sending**: Client calls `SendMessage` RPC method
4. **Receiving**: Background goroutine polls `ReceiveMessage` every second
5. **Concurrency**: Mutex ensures thread-safe message storage


## Usage

### Start the Server
```bash
cd server
go run server.go
```

### Run Clients (multiple terminals)
```bash
cd client
go run client.go
```

## Project Structure
```
chat-rpc/
├── api/
│   └── api.go          # RPC interface definitions
├── server/
│   └── server.go       # RPC server implementation
├── client/
│   └── client.go       # RPC client implementation
├── go.mod              # Go module file
└── README.md
```

## Use Cases

- **Learning RPC concepts** and distributed systems
- **Understanding Go concurrency** with goroutines and mutexes
- **Building real-time applications** without WebSockets
- **Prototyping chat systems** for larger applications
- **Teaching client-server architecture** fundamentals

## Further Enhancements

### Core Features
- Message persistence (SQLite/PostgreSQL)
- User authentication and sessions
- Chat rooms/channels
- Private messaging
- Message history for new users
