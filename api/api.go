package api

import "time"

type Message struct {
	Username  string
	Text      string
	Timestamp time.Time
}

type SendArgs struct {
	Username string
	Text     string
}

type GetArgs struct {
	FromIndex int
}

type ChatService interface {
	SendMessage(args *SendArgs, reply *bool) error
	RecieveMessage(args *GetArgs, reply *[]Message) error
}
