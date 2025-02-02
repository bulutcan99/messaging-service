package entity

import (
	"time"

	"github.com/google/uuid"
)

type StudentRole string

type User struct {
	ID                    uuid.UUID   `json:"id"`
	BoardID               string      `json:"boardID"`
	Username              string      `json:"username"`
	Email                 string      `json:"email"`
	IsEmailVerified       bool        `json:"isEmailVerified"`
	SchoolEmail           string      `json:"schoolEmail"`
	IsSchoolEmailVerified bool        `json:"isSchoolEmailVerified"`
	PhoneNumber           string      `json:"phoneNumber"`
	IsPhoneNumberVerified bool        `json:"isPhoneNumberVerified"`
	Password              string      `json:"password"`
	Avatar                string      `json:"avatar,omitempty"`
	Role                  StudentRole `json:"role"`
	CreatedAt             time.Time   `json:"createdAt"`
	UpdatedAt             time.Time   `json:"updatedAt"`
	DeletedAt             *time.Time  `json:"deletedAt,omitempty"`
}
