package dto

import (
	"time"
)

type ClientRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoomRes struct {
	ID    string      `json:"id"`
	Name  string      `json:"name"`
	Users []ClientRes `json:"users"`
}

type CreateMassageReq struct {
	Text    string
	UserID  string
	RoomID  string
	CreatAt time.Time
}

type CreateMassageRes struct {
	Text    string
	UserID  string
	RoomID  string
	CreatAt time.Time
}
