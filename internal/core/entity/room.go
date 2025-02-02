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

func (r *Room) RemoveClientFromRoom(client *Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.Clients[client.ID]; !exists {
		return fmt.Errorf("client %s is not in room %s", client.ID, r.ID)
	}

	delete(r.Clients, client.ID)

	client.RoomID = ""
	client.Status = Disconnected
	if client.Message != nil {
		close(client.Message)
		client.Message = nil
	}

	zap.S().Debugf("Client %s removed from room %s", client.ID, r.ID)
	return nil
}

func (r *Room) Run(hub *Hub) {
	r.wg.Add(1)
	defer r.wg.Done()

	for {
		select {
		case client := <-r.Register:
			if _, ok := r.Clients[client.ID]; !ok {
				err := r.AddClientToRoom(client)
				if err != nil {
					zap.S().Errorf("error adding client %s to room %s: %v", client.ID, r.ID, err)
				}
				hub.BroadcastRoomUpdate(fmt.Sprintf("Client %s has joined the room", client.ID))
			}

		case client := <-r.Unregister:
			if _, ok := r.Clients[client.ID]; ok {
				err := r.RemoveClientFromRoom(client)
				if err != nil {
					zap.S().Errorf("error removing client %s from room %s: %v", client.ID, r.ID, err)
				}

				zap.S().Debugf("Client %s left room %s", client.ID, r.ID)

				r.mu.Lock()
				delete(r.Clients, client.ID)
				clientCount := len(r.Clients)
				r.mu.Unlock()

				if clientCount == 0 {
					zap.S().Debugf("Room %s is empty, closing", r.ID)
					hub.RoomUnregister <- r.ID
					return
				}
			}
		}
	}
}
