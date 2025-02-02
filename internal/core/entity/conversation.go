package entity

import "github.com/google/uuid"

type Conversation struct {
	ID               uuid.UUID `json:"id"`
	ConversationType string    `json:"conversation_type"` // direct, group
	Participants     []string  `json:"participants"`
	Admins           []string  `json:"admins"`
	GroupName        string    `json:"group_name"`
}
