package chat

import (
	"fmt"
	"go.uber.org/zap"
	"sync"
)

const MaxClientsInRoom = 2

type Room struct {
	ID         string       `json:"id"`
	Name       string       `json:"name"`
	Register   chan *Client `json:"-"`
	Unregister chan *Client `json:"-"`
	Broadcast  chan *ListenerMessage
	Clients    map[string]*Client `json:"clients"`
	mu         sync.Mutex
	wg         sync.WaitGroup
}

func NewRoom(id, name string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *ListenerMessage, 1000),
		Clients:    make(map[string]*Client),
	}
}

func (r *Room) Copy() Room {
	return Room{
		ID:      r.ID,
		Name:    r.Name,
		Clients: r.Clients,
	}
}

func (r *Room) AddClientToRoom(client *Client) error {
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

func (r *Room) AddAdminToRoom(admin *Client) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	if _, exists := r.Clients[admin.ID]; exists {
		return fmt.Errorf("admin %s is already in room %s", admin.ID, r.ID)
	}

	r.Clients[admin.ID] = admin
	admin.RoomID = r.ID
	admin.Status = Connected

	zap.S().Debugf("Admin %s added to room %s", admin.ID, r.ID)
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
				if client.Role == Admin {
					err := r.AddAdminToRoom(client)
					if err != nil {
						zap.S().Errorf("error adding admin %s to room %s: %v", client.ID, r.ID, err)
					}
					hub.BroadcastRoomUpdate(fmt.Sprintf("Admin %s has joined the room", client.ID))
				} else {
					err := r.AddClientToRoom(client)
					if err != nil {
						zap.S().Errorf("error adding client %s to room %s: %v", client.ID, r.ID, err)
					}
					hub.BroadcastRoomUpdate(fmt.Sprintf("Client %s has joined the room", client.ID))
				}
			}

		case client := <-r.Unregister:
			if _, ok := r.Clients[client.ID]; ok {
				err := r.RemoveClientFromRoom(client)
				if err != nil {
					zap.S().Errorf("error removing client %s from room %s: %v", client.ID, r.ID, err)
				}

				if client.Role == User {
					zap.S().Debugf("User %s left room %s, closing room", client.ID, r.ID)

					r.mu.Lock()
					for _, admin := range r.Clients {
						if admin.Role == Admin {
							zap.S().Debugf("Admin %s removed from room %s", admin.ID, r.ID)
							if admin.Message != nil {
								close(admin.Message)
								admin.Message = nil
							}
							delete(r.Clients, admin.ID)
						}
					}
					r.mu.Unlock()

					hub.RoomUnregister <- r.ID
					return
				} else if client.Role == Admin {
					zap.S().Debugf("Admin %s left room %s", client.ID, r.ID)

					r.mu.Lock()
					clientCount := len(r.Clients)
					r.mu.Unlock()

					if clientCount == 0 {
						zap.S().Debugf("Room %s is empty, closing", r.ID)
						hub.RoomUnregister <- r.ID
						return
					}
				}
			}

		case message := <-r.Broadcast:
			r.mu.Lock()
			for _, client := range r.Clients {
				select {
				case client.Message <- message:
				default:
					zap.S().Errorf("Failed to send message to client %s: message channel full", client.ID)
					close(client.Message)
					delete(r.Clients, client.ID)
				}
			}
			r.mu.Unlock()
		}
	}
}
