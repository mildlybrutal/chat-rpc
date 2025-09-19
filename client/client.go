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

func main() {
	fmt.Print("Enter your username: ")
	username, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	username = strings.TrimSpace(username)

	client, err := rpc.Dial("tcp", "localhost:8080")

	if err != nil {
		fmt.Println("Error connecting: ", err)
		return
	}

	defer client.Close()
	//goroutine to recieve messages  - runs in background
	go func() {
		var lastIndex int
		for {
			args := &api.GetArgs{FromIndex: lastIndex}
			var newMessages []api.Message
			err = client.Call("ChatServer.ReceiveMessage ", args, &newMessages)
			if err != nil {
				fmt.Println("Error recieving messages: ", err)
			} else {
				for _, msg := range newMessages {
					fmt.Printf("[%s] %s: %s\n", msg.Timestamp.Format("15:04"), msg.Username, msg.Text)
				}
				if len(newMessages) > 0 {
					lastIndex += len(newMessages)
				}

			}
			time.Sleep(1 * time.Second)
		}
	}()

	// Main goroutine for sending messages - reads user input from console and sends messages to the server via SendMessage RPC method
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("type your message")
	for scanner.Scan() {
		messageText := scanner.Text()
		if len(messageText) == 0 {
			continue
		}

		args := &api.SendArgs{
			Username: username,
			Text:     messageText,
		}
		var reply bool // The server will send a boolean to confirm
		err := client.Call("ChatServer.SendMessage", args, &reply)

		if err != nil {
			fmt.Println("RPC call failed: ", err)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Reading input failed:", err)
	}
}
