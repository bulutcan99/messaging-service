package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type BlockList struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	BlockerID bson.ObjectID `bson:"blocker_id"`
	BlockedID bson.ObjectID `bson:"blocked_id"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty"`
}
