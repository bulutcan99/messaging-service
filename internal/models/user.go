package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type StudentRole string

const (
	RoleStudent StudentRole = "STUDENT"
)

type User struct {
	ID                    bson.ObjectID `bson:"_id,omitempty"`
	BoardID               string        `bson:"board_id"`
	Username              string        `bson:"username"`
	Email                 string        `bson:"email"`
	IsEmailVerified       bool          `bson:"is_email_verified"`
	SchoolEmail           *string       `bson:"school_email,omitempty"`
	IsSchoolEmailVerified bool          `bson:"is_school_email_verified"`
	PhoneNumber           *string       `bson:"phone_number,omitempty"`
	IsPhoneNumberVerified bool          `bson:"is_phone_number_verified"`
	Password              string        `bson:"password"`
	Avatar                *string       `bson:"avatar,omitempty"`
	Role                  StudentRole   `bson:"role"`
	LastSeenAt            *time.Time    `bson:"last_seen_at,omitempty"`
	CreatedAt             time.Time     `bson:"created_at"`
	UpdatedAt             time.Time     `bson:"updated_at"`
	DeletedAt             *time.Time    `bson:"deleted_at,omitempty"`
}
