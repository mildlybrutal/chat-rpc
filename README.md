# RPC Chat Application

A real-time chat application built with Go's native RPC package, demonstrating distributed systems concepts with concurrent client-server architecture and multi-room chat functionality.


## Overview

This project implements a feature-rich chat system using Go's built-in RPC (Remote Procedure Call) framework. Multiple clients can connect to a central server to send and receive messages in real-time across different chat rooms.

## Architecture

```
┌─────────────┐    RPC Calls    ┌─────────────┐
│   Client    │ ◄─────────────► │   Server    │
│   ┌─────────┤                 ├─────────┐   │
│   │ Sender  │                 │ Room    │   │
│   │ Thread  │                 │ Manager │   │
│   ├─────────┤                 ├─────────┤   │
│   │Receiver │                 │ Message │   │
│   │ Thread  │                 │ Store   │   │
│   ├─────────┤                 ├─────────┤   │
│   │Commands │                 │ User    │   │
│   │(/join   │                 │ Tracker │   │
│   │/leave)  │                 │         │   │
│   └─────────┘                 └─────────┘   │
└─────────────┘                 └─────────────┘
```


### Components:
- **API Layer**: Shared interface definitions and data structures
- **Server**: RPC service handling message storage, room management, and user tracking
- **Client**: Concurrent sender/receiver with command-based user interface
- **Room System**: Multi-room chat with dynamic room creation and user management

## How It Works

1. **Server** starts and listens on port 8080
2. **Clients** connect via TCP and register with username  
3. **Default Room**: All users start in the "general" room
4. **Room Management**: Users can create/join rooms dynamically using commands
5. **Sending**: Client calls `SendMessage` RPC method to current room
6. **Receiving**: Background goroutine polls `ReceiveMessage` every second per room
7. **Concurrency**: Mutex ensures thread-safe message storage and user tracking

## Features

### **Multi-Room Chat System**
- **Dynamic Room Creation**: Rooms are created automatically when first joined
- **Default Room**: All users start in the `#general` room
- **Room Switching**: Seamlessly switch between multiple rooms
- **Per-Room Message History**: Each room maintains its own message history
- **User Tracking**: Server tracks which users are in which rooms

### **Interactive Commands**
- `/join <room>` - Join or create a new chat room
- `/leave` - Leave current room (returns to #general)
- `/rooms` - List all available rooms with current room indicator
- `/users` - Show users in current room (placeholder for future enhancement)
- `/help` - Display all available commands

###**Real-Time Messaging**
- Instant message delivery within rooms
- Message timestamps and room indicators
- Concurrent message handling across multiple rooms


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
│   └── api.go          # RPC interface definitions & room structures
├── server/
│   └── server.go       # RPC server with room management
├── client/
│   └── client.go       # RPC client with command interface
├── go.mod              # Go module file
└── README.md
```

## Technical Implementation

### Server-Side Room Management
- **Room Storage**: `map[string][]api.Message` - Messages organized by room
- **User Tracking**: Bidirectional mapping of users to rooms
- **Thread Safety**: Mutex locks protect concurrent room operations
- **Dynamic Creation**: Rooms created on-demand when users join

### Client-Side Features
- **Command Parser**: Handles `/` prefixed commands vs. regular messages
- **Room State**: Tracks current room and message indices per room
- **Real-time Updates**: Polls for new messages every second in current room
- **User Interface**: Clean room indicators and status messages

## Use Cases

- **Learning RPC concepts** and distributed systems
- **Understanding Go concurrency** with goroutines and mutexes
- **Building real-time applications** without WebSockets
- **Prototyping chat systems** for larger applications
- **Teaching client-server architecture** fundamentals
- **Demonstrating room-based chat systems** like Discord/Slack
- **Multi-user collaboration** and real-time communication

## Further Enhancements

### Implemented
- **Chat rooms/channels** - Multi-room support with dynamic creation
- **Room management** - Join, leave, list rooms functionality
- **User tracking** - Server maintains user-room relationships
- **Command interface** - Slash commands for room operations

### Potential Improvements 
- Message persistence
- User authentication and sessions
- User list display per room
- Room moderation and permissions

