package entity

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type Room struct {
	ID         string           `json:"id"`
	Name       string           `json:"name"`
	Clients    map[string]*User `json:"clients"`
	Register   chan *User       `json:"-"`
	Unregister chan *User       `json:"-"`
	mu         sync.Mutex
	wg         sync.WaitGroup
}

func NewRoom(id, name string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		Clients:    make(map[string]*User),
		Register:   make(chan *User),
		Unregister: make(chan *User),
	}
}

func (r *Room) Copy() Room {
	return Room{
		ID:      r.ID,
		Name:    r.Name,
		Clients: r.Clients,
	}
}

func (r *Room) AddClientToRoom(client *User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.Clients[client.ID]; exists {
		return fmt.Errorf("client %s is already in room %s", client.ID, r.ID)
	}

	r.Clients[client.ID] = client
	client.RoomID = r.ID
	client.Status = Connected

	zap.S().Debugf("Client %s added to room %s", client.ID, r.ID)
	return nil
}
