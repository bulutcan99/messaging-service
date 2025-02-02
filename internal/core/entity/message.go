package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	Message        string    `json:"message"`
	SenderUserID   string    `json:"sender_user_id"`
	File           string    `json:"file"`
	FileName       string    `json:"file_name"`
	FileType       string    `json:"file_type"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
