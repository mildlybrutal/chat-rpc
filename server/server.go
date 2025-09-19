package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"

	"github.com/chat-rpc/api"
)

// main rpc service
type ChatServer struct {
	messages []api.Message
	mu       sync.Mutex
}

// RPC method to receive and store a new chat message.
func (c *ChatServer) SendMessage(args *api.SendArgs, reply *bool) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	newMessage := api.Message{
		Username:  args.Username,
		Text:      args.Text,
		Timestamp: time.Now(),
	}

	c.messages = append(c.messages, newMessage)
	log.Printf("Received message from %s: %s\n", args.Username, args.Text)
	*reply = true
	return nil
}

// RPC method to retrieve new messages from the chat history.
func (c *ChatServer) RecieveMessage(args *api.GetArgs, reply *[]api.Message) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	if args.FromIndex < 0 || args.FromIndex > len(c.messages) {
		return nil
	}

	*reply = c.messages[args.FromIndex:]
	return nil
}

func main() {
	chat := new(ChatServer)
	rpc.Register(chat)

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Print("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting conn:", err)
			continue
		}

		go rpc.ServeConn(conn) //concurreny
	}

}
