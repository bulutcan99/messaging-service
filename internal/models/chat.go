package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatType string

const (
	Individual ChatType = "INDIVIDUAL"
	Group      ChatType = "GROUP"
)

type Participant struct {
	UserID   bson.ObjectID `bson:"user_id"`
	JoinedAt time.Time     `bson:"joined_at"`
	LeftAt   *time.Time    `bson:"left_at,omitempty"`
	IsAdmin  bool          `bson:"is_admin"`
}

type Chat struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	Type         ChatType      `bson:"type"`
	Name         *string       `bson:"name,omitempty"`        // For group chats
	Description  *string       `bson:"description,omitempty"` // For group chats
	Participants []Participant `bson:"participants"`
	CreatedAt    time.Time     `bson:"created_at"`
	UpdatedAt    time.Time     `bson:"updated_at"`
	DeletedAt    *time.Time    `bson:"deleted_at,omitempty"`
	LastMessage  *Message      `bson:"last_message,omitempty"`
	IsActive     bool          `bson:"is_active"`
}
