package message

import (
	"go_grpc_realtime/lib/core/grpcgen"
	"time"
)

type MessageRoomQuery struct {
	ID         uint   `json:"id"`
	RoomName   string `json:"room_name"`
	IsOneToOne bool   `json:"is_one_to_one"`
}

type MessageQuery struct {
	MessageId   uint      `json:"message_id"`
	CreatedTime time.Time `json:"created_time"`
	Message     string    `json:"message"`
	RoomId      uint      `json:"room_id"`
	SenderId    uint      `json:"sender_id"`
	SenderName  string    `json:"sender_name"`
	SenderEmail string    `json:"sender_email"`
}

type MessageListener struct {
	UserId  uint
	Channel chan *grpcgen.Message
}
