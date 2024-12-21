package chat

import (
	"time"
)

type ListenerMessage struct {
	Type     ListenerEvent `json:"type"`
	Data     any           `json:"data"`
	IsClient bool          `json:"is_client"`
}

type Message struct {
	ID             string    `json:"id,omitempty"`
	RoomID         string    `json:"room_id"`
	TextEvent      TextEvent `json:"event,omitempty"`
	Text           string    `json:"text,omitempty"`
	SenderUserID   string    `json:"sender_user_id,omitempty"`
	ReceiverUserID string    `json:"receiver_user_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ClientMessage struct {
	RoomID    string `json:"room_id"`
	TextEvent int    `json:"text_event"`
	Text      string `json:"text"`
}

type ClientInAndOutRoom struct {
	ID string `json:"id"`
}
