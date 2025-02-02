package entity

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ParticipantID uuid.UUID `json:"id"`
	Timestamp     time.Time `json:"timestamp"`
}

type MessageEvent struct {
	ID        uuid.UUID `json:"id"`
	MessageID uuid.UUID `json:"message_id"`
	Delivered []Event   `json:"delivered"`
	Read      []Event   `json:"read"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
