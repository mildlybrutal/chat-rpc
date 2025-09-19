package api

import "time"

type Message struct {
	UserName  string
	Text      string
	Timestamp time.Time
}

type SendArgs struct {
	UserName string
	Text     string
}

type GetArgs struct {
	FromIndex int
}

type ChatService interface {
	SendMessage(args *SendArgs, reply *bool) error
	RecieveMessage(args *GetArgs, reply *bool) error
}
