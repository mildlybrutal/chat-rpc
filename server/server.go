package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"slices"
	"sync"
	"time"

	"github.com/chat-rpc/api"
)

// main rpc service
type ChatServer struct {
	messages  []api.Message
	rooms     map[string][]api.Message
	userRooms map[string][]string // room which user is in
	roomUsers map[string][]string //which users are in room
	mu        sync.Mutex
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
func (c *ChatServer) ReceiveMessage(args *api.GetArgs, reply *[]api.Message) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	if args.FromIndex < 0 || args.FromIndex > len(c.messages) {
		return nil
	}

	*reply = c.messages[args.FromIndex:]
	return nil
}

func (c *ChatServer) ListRooms(args *struct{}, reply *api.ListRoomsReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	rooms := make([]api.RoomInfo, 0, len(c.rooms))

	for name := range c.rooms {
		rooms = append(rooms, api.RoomInfo{RoomName: name})
	}

	reply.Rooms = rooms

	return nil
}

func (c *ChatServer) JoinRoom(args *api.RoomInfo, reply *api.JoinRoomReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	//initializes room if not exists in first place
	if _, roomExists := c.rooms[args.RoomName]; !roomExists {
		c.rooms[args.RoomName] = make([]api.Message, 0)
		c.roomUsers[args.RoomName] = make([]string, 0)
	}

	//prevents duplicateds , checks if user is already in room
	if slices.Contains(c.roomUsers[args.RoomName], args.Username) {
		return nil // Already in room
	}

	//bi-directional
	//add user to room
	c.roomUsers[args.RoomName] = append(c.roomUsers[args.RoomName], args.Username)
	//add room to user's list
	c.userRooms[args.Username] = append(c.userRooms[args.Username], args.RoomName)

	*reply = api.JoinRoomReply{
		Success:   true,
		Message:   fmt.Sprintf("joined room %s", args.RoomName),
		UserCount: len(c.roomUsers[args.RoomName]),
	}
	//fmt.Printf("User %s joined room %s", args.Username, args.RoomName)
	return nil
}

func (c *ChatServer) LeaveRoom(args *api.RoomInfo, reply *api.LeaveRoomReply) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// removes user from room
	c.roomUsers[args.RoomName] = slices.DeleteFunc(c.roomUsers[args.RoomName],
		func(user string) bool { return user == args.Username })
	// remove room from user
	c.userRooms[args.Username] = slices.DeleteFunc(c.userRooms[args.Username],
		func(room string) bool { return room == args.RoomName })

	*reply = api.LeaveRoomReply{
		Success:        true,
		Message:        fmt.Sprintf("Left room %s", args.RoomName),
		RemainingUsers: len(c.roomUsers[args.RoomName]),
	}
	return nil
}

func main() {
	chat := &ChatServer{
		rooms:     make(map[string][]api.Message),
		userRooms: make(map[string][]string),
		roomUsers: make(map[string][]string),
	}
	rpc.Register(chat)

	listener, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting conn:", err)
			continue
		}

		go rpc.ServeConn(conn) //concurreny
	}

}
