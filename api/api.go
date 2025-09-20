package api

import "time"

type Message struct {
	Username  string
	Text      string
	RoomName  string
	Timestamp time.Time
}

type SendArgs struct {
	Username string
	Text     string
	RoomName string
}

type GetArgs struct {
	FromIndex int
	RoomName  string
}

type RoomInfo struct {
	RoomName string
	Username string
}

type ListRoomsReply struct {
	Rooms []RoomInfo
}

type JoinRoomReply struct {
	Success   bool
	Message   string
	UserCount int
}

type LeaveRoomReply struct {
	Success        bool
	Message        string
	RemainingUsers int
}

type ChatService interface {
	SendMessage(args *SendArgs, reply *bool) error
	ReceiveMessage(args *GetArgs, reply *[]Message) error
	ListRooms(args *struct{}, reply *ListRoomsReply) error
	JoinRoom(args *RoomInfo, reply *JoinRoomReply) error
	LeaveRoom(args *RoomInfo, reply *LeaveRoomReply) error
}
