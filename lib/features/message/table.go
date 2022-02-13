package message

import (
	"go_grpc_realtime/lib/features/user"
	"time"
)

type RoomTbl struct {
	ID          uint       `json:"id" gorm:"primary_key"`
	Name        string     `json:"name" gorm:"not null"`
	IsOneToOne  bool       `json:"is_one_to_one" gorm:"not null"`
	LastUpdated *time.Time `json:"last_update"`
}

type RoomMembersTbl struct {
	RoomId uint `json:"room_id" gorm:"primaryKey;autoIncrement:false"`
	UserId uint `json:"user_id" gorm:"primaryKey;autoIncrement:false"`
}

type MessageTbl struct {
	ID        uint           `json:"id" gorm:"primary_key"`
	RoomId    uint           `json:"room_id" gorm:"not null"`
	Room      RoomMembersTbl `json:"room" gorm:"foreignKey:RoomId"`
	SenderId  uint           `json:"sender_id" gorm:"not null"`
	Sender    user.UserTbl   `json:"sender" gorm:"foreignKey:SenderId"`
	CreatedAt time.Time      `json:"created_at" gorm:"default:current_timestamp"`
	Body      string         `json:"body" gorm:"not null"`
}
