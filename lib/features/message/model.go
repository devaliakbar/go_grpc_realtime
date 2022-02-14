package message

type MessageRoomQuery struct {
	ID         uint   `json:"id"`
	RoomName   string `json:"room_name"`
	IsOneToOne bool   `json:"is_one_to_one"`
}
