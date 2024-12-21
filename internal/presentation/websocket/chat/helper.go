package chat

import "fmt"

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
