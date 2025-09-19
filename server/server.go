package main

import (
	"time"

	"github.com/chat-rpc/api"
)

type ChatServer struct{
	messages []Message
}

func (c *ChatServer) SendMessage(api.SendArgs,reply *bool) error{
	newMessage := api.Message{
		UserName: args.Username,
		Text: args.Text,
		Timestamp:time.Now() ,
	}

	c.messages = append(c.messages, newMessage)
	log.Printf("Received message from %s: %s\n", args.Username, args.Text)
	*reply=true
	return nil
}