package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"strings"
	"time"

	"github.com/chat-rpc/api"
)

var (
	currentRoom   = "general"
	roomLastIndex = make(map[string]int)
	client        *rpc.Client
	username      string
)

func main() {
	fmt.Print("Enter your username: ")
	reader := bufio.NewReader(os.Stdin)
	username, _ = reader.ReadString('\n')
	username = strings.TrimSpace(username)

	client, err := rpc.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Error connecting: ", err)
		return
	}

	defer client.Close()

	joinRoom("general") //default when user joins
	listRooms()

	//goroutine to recieve messages  - runs in background
	go messageReceiver()

	// Main goroutine for sending messages - reads user input from console and sends messages to the server via SendMessage RPC method
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("[#%s] Type your message (or /help for commands):\n", currentRoom)
	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())
		if len(input) == 0 {
			continue
		}

		if strings.HasPrefix(input, "/") {
			handleCommand(input)
		} else {
			SendMessage(input)
		}

	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Reading input failed:", err)
	}
}

func messageReceiver() {
	for {
		args := &api.GetArgs{
			FromIndex: roomLastIndex[currentRoom],
			RoomName:  currentRoom,
		}

		var newMessages []api.Message
		err := client.Call("ChatServer.ReceiveMessage ", args, &newMessages)

		if err != nil {
			fmt.Println("Error recieving messages: ", err)
		} else {
			for _, msg := range newMessages {
				fmt.Printf("[#%s][%s] %s: %s\n", msg.RoomName, msg.Timestamp.Format("15:04"), msg.Username, msg.Text)
			}
			if len(newMessages) > 0 {
				roomLastIndex[currentRoom] += len(newMessages)
			}

		}
		time.Sleep(1 * time.Second)
	}
}

func SendMessage(text string) {
	args := &api.SendArgs{
		Username: username,
		Text:     text,
		RoomName: currentRoom,
	}
	var reply bool // The server will send a boolean to confirm
	err := client.Call("ChatServer.SendMessage", args, &reply)

	if err != nil {
		fmt.Println("failed to send message: ", err)
	} else if reply {
		fmt.Printf("[#%s] You: %s\n", currentRoom, text)
	}
}

func handleCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]

	switch command {
	case "/help":
		showHelp()
	case "/join":
		if len(parts) < 2 {
			fmt.Println("Usage: /join <roomname>")
			return
		}
		joinRoom(parts[1])
	case "/leave":
		leaveRoom()
	case "/rooms":
		listRooms()
	case "/users":
		showUsers()
	default:
		fmt.Printf("Unknown command: %s (type /help for available commands)\n", command)
	}

}

func joinRoom(roomName string) {
	args := &api.RoomInfo{
		Username: username,
		RoomName: roomName,
	}

	var reply api.JoinRoomReply

	err := client.Call("ChatServer.JoinRoom", args, &reply)

	if err != nil {
		fmt.Printf("Failed to join room %s: %v\n", roomName, err)
		return
	}

	if reply.Success {
		oldRoom := currentRoom
		currentRoom := roomName

		if _, exists := roomLastIndex[currentRoom]; !exists {
			roomLastIndex[currentRoom] = 0
		}
		fmt.Printf("Joined room #%s (%d users online)\n", roomName, reply.UserCount)
		if oldRoom != roomName {
			fmt.Printf("Switched from #%s to #%s\n", oldRoom, roomName)
		}
	} else {
		fmt.Printf("Failed to join room %s: %s\n", roomName, reply.Message)
	}
}

func leaveRoom() {
	if currentRoom == "general" {
		fmt.Println("Cannot leave the general room")
		return
	}

	args := &api.RoomInfo{
		Username: username,
		RoomName: currentRoom,
	}

	var reply api.LeaveRoomReply

	err := client.Call("ChatServer.LeaveRoom", args, &reply)

	if err != nil {
		fmt.Printf("Failed to leave room: %v\n", err)
		return
	}

	if reply.Success {
		fmt.Printf("Left room #%s", currentRoom)
		currentRoom = "general"
		fmt.Printf("Switched to #%s\n", currentRoom)
	} else {
		fmt.Printf("Failed to leave room: %s\n", reply.Message)
	}
}

func listRooms() {
	args := &struct{}{}

	var reply api.ListRoomsReply
	err := client.Call("ChatServer.ListRooms", args, &reply)

	if err != nil {
		fmt.Println("Failed to list rooms:", err)
		return
	}

	fmt.Println("Available Rooms:")
	if len(reply.Rooms) == 0 {
		fmt.Println("no rooms available")
	} else {
		for _, room := range reply.Rooms {
			marker := "	"
			if room.RoomName == currentRoom {
				marker = "* "
			}
			fmt.Printf("%s#%s\n", marker, room.RoomName)
		}
	}
	fmt.Printf("Currently in: #%s\n\n", currentRoom)
}

func showUsers() {
	fmt.Printf("Users in #%s:\n", currentRoom)
	fmt.Println("  (User list feature not implemented yet)")
}

func showHelp() {
	fmt.Println("\nAvailable Commands:")
	fmt.Println("  /join <room>  - Join or create a room")
	fmt.Println("  /leave        - Leave current room (return to #general)")
	fmt.Println("  /rooms        - List all available rooms")
	fmt.Println("  /users        - Show users in current room")
	fmt.Println("  /help         - Show this help")
	fmt.Println("  <message>     - Send message to current room")
	fmt.Printf("\nCurrently in: #%s\n\n", currentRoom)
}
