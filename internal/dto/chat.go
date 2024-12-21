package dto

import "time"

type ChatDTO struct {
	ID           string      `json:"id"`
	Type         string      `json:"type"`
	Name         *string     `json:"name,omitempty"`
	Description  *string     `json:"description,omitempty"`
	Participants []string    `json:"participants"`
	LastMessage  *MessageDTO `json:"last_message,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}
