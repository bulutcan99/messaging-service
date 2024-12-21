package chat

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"sync"
)

//todo: kullanicilarin disconnect kisimlari duzeltilecek (normal userlarin tum her seyleri kapatilirken
// adminlerin sadece odadan cikmansi gerekli)

type RoomResponse struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Clients []string `json:"clients"`
}

type RoomUpdate struct {
	Message string         `json:"message"`
	Rooms   []RoomResponse `json:"rooms"`
}

type Hub struct {
	Rooms          sync.Map
	Register       chan *Client
	Unregister     chan *Client
	RoomRegister   chan *Client
	RoomUnregister chan string
	Broadcast      chan *ListenerMessage
	AdminClients   map[string]*Client
	mu             sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		RoomRegister:   make(chan *Client),
		RoomUnregister: make(chan string),
		Broadcast:      make(chan *ListenerMessage, 1000),
		AdminClients:   make(map[string]*Client),
	}
}

func (h *Hub) RegisterAdmin(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.AdminClients[client.ID] = client
}

func (h *Hub) UnregisterAdmin(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.AdminClients, client.ID)
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

func (h *Hub) CreateRoom(cl *Client) error {
	if cl.Role != User {
		return errors.New("only users can create rooms")
	}

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
