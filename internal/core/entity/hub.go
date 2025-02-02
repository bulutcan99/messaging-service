package entity

import (
	"errors"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type Hub struct {
	Rooms sync.Map
	mu    sync.Mutex
}

func NewHub() *Hub {
	return &Hub{}
}

func (h *Hub) CreateRoom(cl *User) error {
	name := fmt.Sprintf("%s's Room", cl.Name)
	if ok := h.IsValidRoom(cl.ID, name); !ok {
		return errors.New("room already exists")
	}

	room := NewRoom(cl.ID, name)
	h.Rooms.Store(cl.ID, room)
	h.BroadcastRoomUpdate(fmt.Sprintf("Room created: %s", cl.ID))

	go room.Run(h)
	room.Register <- cl
	return nil
}

func (h *Hub) BroadcastRoomUpdate(message string) {
	rooms := make([]RoomResponse, 0)

	h.Rooms.Range(func(_, value interface{}) bool {
		room, ok := value.(*Room)
		if !ok {
			zap.S().Error("Failed to cast value to Room")
			return false
		}

		cls := make([]string, 0)
		for _, cl := range room.Clients {
			cls = append(cls, cl.ID)
		}

		toRoomResponse := RoomResponse{
			ID:      room.ID,
			Name:    room.Name,
			Clients: cls,
		}

		rooms = append(rooms, toRoomResponse)
		return true
	})

	update := ListenerMessage{
		Type: EventHubUpdate,
		Data: RoomUpdate{
			Message: message,
			Rooms:   rooms,
		},
		IsClient: false,
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	for _, client := range h.AdminClients {
		if client.Message != nil {
			select {
			case client.Message <- &update:
			default:
				zap.S().Errorf("Failed to send message to admin %s: message channel full", client.ID)
			}
		} else {
			zap.S().Errorf("Admin %s does not have a message channel", client.ID)
		}
	}
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

func (h *Hub) IsAvailableClient(roomID, clientID string) *User {
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

func (h *Hub) DeleteRoom(roomID string) error {
	h.Rooms.Delete(roomID)
	h.BroadcastRoomUpdate(fmt.Sprintf("Room deleted: %s", roomID))
	return nil
}

func (h *Hub) Run() {
	for {
		select {
		case cl := <-h.Register:
			if cl.Role == Admin {
				h.RegisterAdmin(cl)
				zap.S().Debugf("Admin %s has joined the hub", cl.ID)
				h.BroadcastRoomUpdate(fmt.Sprintf("Admin %s has joined the hub.", cl.ID))
			} else {
				zap.S().Errorf("Only admins can register in the hub.")
			}

		case cl := <-h.RoomRegister:
			if err := h.CreateRoom(cl); err != nil {
				zap.S().Errorf("Failed to create room: %v", err)
			}

		case roomID := <-h.RoomUnregister:
			if err := h.DeleteRoom(roomID); err != nil {
				zap.S().Errorf("Failed to delete room: %v", err)
			}
		}
	}
}
