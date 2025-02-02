package entity

import (
	"errors"
	"fmt"
	"sync"
)

type Hub struct {
	Rooms sync.Map
	mu    sync.Mutex
}

func NewHub() *Hub {
	return &Hub{}
}

func (h *Hub) CreateRoom(user *User) error {
	name := fmt.Sprintf("%s's Room", user.Name)
	if ok := h.IsValidRoom(user.ID, name); !ok {
		return errors.New("room already exists")
	}

	room := NewRoom(user.ID, name)
	h.Rooms.Store(user.ID, room)
	return nil
}

func (h *Hub) IsValidRoom(roomID, name string) bool {
	valid := true
	h.Rooms.Range(func(_, value interface{}) bool {
		room := value.(*Room)
		if room.ID == roomID || room.Name == name {
			valid = false
			return false
		}
		return true
	})
	return valid
}

func (h *Hub) IsCreatedRoom(roomID string) bool {
	_, ok := h.Rooms.Load(roomID)
	return ok
}

func (h *Hub) IsAvailableClient(roomID, clientID string) *Client {
	if roomInterface, ok := h.Rooms.Load(roomID); ok {
		room := roomInterface.(*Room)
		for _, c := range room.Clients {
			if c.ID != clientID {
				return c
			}
		}
	}
	return nil
}

func (h *Hub) IsAvailableClientID(roomID, clientID string) (string, error) {
	if roomInterface, ok := h.Rooms.Load(roomID); ok {
		room := roomInterface.(*Room)
		for _, c := range room.Clients {
			if c.ID != clientID {
				return c.ID, nil
			}
		}
	}
	return "", fmt.Errorf("no available client %s, in this room %s", clientID, roomID)
}
