package entity

type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status Status `json:"status,omitempty"`
	RoomID string `json:"room_id,omitempty"`
}

func NewUser(id, name string) *User {
	return &User{
		ID:     id,
		Name:   name,
		Status: Idle,
		RoomID: "",
	}
}
