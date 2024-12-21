package chat

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"time"
)

type Role int

const (
	Admin Role = iota
	User
)

type Client struct {
	Conn    *websocket.Conn       `json:"-"`
	Message chan *ListenerMessage `json:"-"`
	RoomID  string                `json:"room_id"`
	ID      string                `json:"id,omitempty"`
	Name    string                `json:"name,omitempty"`
	Status  Status                `json:"status,omitempty"`
	Role    Role                  `json:"role"`
}

func (c *Client) NotifyRoomFull(roomID string) {
	listenerMsg :=
		&ListenerMessage{
			Type: EventRoomUpdate,
			Data: &Message{
				ID:             uuid.NewString(),
				Text:           "Room is full, cannot join.",
				RoomID:         roomID,
				CreatedAt:      time.Now(),
				ReceiverUserID: c.ID,
			},
			IsClient: false,
		}

	if err := c.Conn.WriteJSON(listenerMsg); err != nil {
		zap.S().Errorf("Error notifying client %s that room is full: %v", c.ID, err)
		c.Conn.Close()
	} else {
		zap.S().Debugf("Notified client %s that room %s is full", c.ID, roomID)
	}
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		if err := c.Conn.WriteJSON(message); err != nil {
			zap.S().Errorf("error writing message: %v", err)
			return
		}
	}
}
func (c *Client) ReadMessage(hub *Hub) {
	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				zap.S().Errorf("Unexpected close error for client %s: %v", c.ID, err)
			} else {
				zap.S().Debugf("Client %s disconnected", c.ID)
			}

			c.Conn.Close()
			break
		}

		var incomeMsg ListenerMessage
		if err := json.Unmarshal(m, &incomeMsg); err != nil {
			zap.S().Errorf("error unmarshaling message from client %s: %v", c.ID, err)
			break
		}

		if !incomeMsg.IsClient {
			zap.S().Debugf("Skipping non-client message for client %s", c.ID)
			continue
		}

		switch incomeMsg.Type {
		case EventRoomRegister:
			var messageData ClientInAndOutRoom
			messageBytes, err := json.Marshal(incomeMsg.Data)
			if err != nil {
				zap.S().Errorf("Failed to marshal incomeMsg.Data: %v", err)
				return
			}

			if err := json.Unmarshal(messageBytes, &messageData); err != nil {
				zap.S().Errorf("Failed to unmarshal AdminRegisterRoom: %v", err)
				return
			}

			zap.S().Debugf("Client %s registering to room %s", c.ID, messageData.ID)

			roomInterface, exists := hub.Rooms.Load(messageData.ID)
			if !exists {
				zap.S().Warnf("Room with ID %s not found", messageData.ID)
				c.NotifyRoomNotFound(messageData.ID)
				return
			}

			room := roomInterface.(*Room)
			if len(room.Clients) >= MaxClientsInRoom {
				zap.S().Infof("Room %s is full, client %s cannot join", messageData.ID, c.ID)
				c.NotifyRoomFull(messageData.ID)
				return
			}

			zap.S().Infof("Registering client %s to room %s", c.ID, messageData.ID)
			room.Register <- c
		case EventRoomUnregister:
			var messageData ClientInAndOutRoom
			messageBytes, err := json.Marshal(incomeMsg.Data)
			if err != nil {
				zap.S().Errorf("Failed to marshal incomeMsg.Data: %v", err)
				return
			}

			if err := json.Unmarshal(messageBytes, &messageData); err != nil {
				zap.S().Errorf("Failed to unmarshal AdminUnregisterRoom: %v", err)
				return
			}

			fmt.Print("MSGDATA: ", messageData)

			zap.S().Debugf("Client %s unregistering from room %s", c.ID, messageData.ID)

			roomInterface, exists := hub.Rooms.Load(messageData.ID)
			if !exists {
				zap.S().Warnf("Room with ID %s not found", messageData.ID)
				c.NotifyRoomNotFound(messageData.ID)
				return
			}

			room := roomInterface.(*Room)

			if _, ok := room.Clients[c.ID]; !ok {
				zap.S().Warnf("Client %s is not in room %s", c.ID, messageData.ID)
				c.NotifyClientNotInRoom(hub, messageData.ID)
				return
			}

			zap.S().Infof("Unregistering client %s from room %s", c.ID, messageData.ID)
			room.Unregister <- c

			c.NotifyUnregisteredFromRoom(hub, messageData.ID)

			if len(room.Clients) == 0 {
				zap.S().Infof("Room %s is empty after client %s left, closing room", messageData.ID, c.ID)
				hub.RoomUnregister <- messageData.ID
			}

		case EventRoomBroadcast:
			var messageData ClientMessage
			messageBytes, err := json.Marshal(incomeMsg.Data)
			if err != nil {
				zap.S().Errorf("Failed to marshal incomeMsg.Data: %v", err)
				return
			}

			if err := json.Unmarshal(messageBytes, &messageData); err != nil {
				zap.S().Errorf("Failed to unmarshal Message: %v", err)
				return
			}

			zap.S().Debugf("Client %s broadcasting message to room %s: %s", c.ID, messageData.RoomID, messageData.Text)
			if roomInterface, exists := hub.Rooms.Load(messageData.RoomID); exists {
				room := roomInterface.(*Room)

				availableClientID, err := hub.IsAvailableClientID(messageData.RoomID, c.ID)
				if err != nil {
					zap.S().Warnf("Room %s has less than 2 clients, message will not be broadcast", messageData.RoomID)
					errMessage := ListenerMessage{
						Type: EventMessageError,
						Data: &Message{
							ID:        uuid.NewString(),
							Text:      "Room must have at least 2 clients to broadcast messages.",
							RoomID:    messageData.RoomID,
							CreatedAt: time.Now(),
						},
						IsClient: false,
					}
					if err := c.Conn.WriteJSON(errMessage); err != nil {
						zap.S().Errorf("Failed to send error message to client %s: %v", c.ID, err)
					}
					return
				}

				serverMsg := &ListenerMessage{
					Type: EventRoomBroadcast,
					Data: &Message{
						ID:             uuid.NewString(),
						Text:           messageData.Text,
						RoomID:         messageData.RoomID,
						TextEvent:      TextEvent(messageData.TextEvent),
						CreatedAt:      time.Now(),
						SenderUserID:   c.ID,
						ReceiverUserID: availableClientID,
					},
					IsClient: true,
				}

				room.Broadcast <- serverMsg
			} else {
				zap.S().Errorf("Room %s not found for broadcasting message", messageData.RoomID)
			}
		default:
			zap.S().Warnf("Client %s sent unknown event type: %d", c.ID, incomeMsg.Type)
		}
	}
}

func (c *Client) NotifyRoomNotFound(roomID string) {
	msg := &ListenerMessage{
		Type: EventRoomUpdate,
		Data: &Message{
			ID:             uuid.NewString(),
			Text:           fmt.Sprintf("Room with ID %s could not be found.", roomID),
			RoomID:         roomID,
			CreatedAt:      time.Now(),
			ReceiverUserID: c.ID,
		},
		IsClient: false,
	}

	if err := c.Conn.WriteJSON(msg); err != nil {
		zap.S().Errorf("Error writing message to client %s: %v", c.ID, err)
	} else {
		zap.S().Infof("Sent RoomNotFound message to client %s for room %s", c.ID, roomID)
	}
}

func (c *Client) NotifyUnregisteredFromRoom(hub *Hub, roomID string) {
	msg := &ListenerMessage{
		Type: EventRoomUpdate,
		Data: &Message{
			ID:             uuid.NewString(),
			Text:           fmt.Sprintf("Client %s disconnected from room %s.", c.ID, roomID),
			RoomID:         roomID,
			CreatedAt:      time.Now(),
			ReceiverUserID: c.ID,
		},
		IsClient: false,
	}

	otherClient := hub.IsAvailableClient(roomID, c.ID)
	if otherClient != nil {
		if err := otherClient.Conn.WriteJSON(msg); err != nil {
			zap.S().Errorf("Error writing RoomNotIn message to client %s: %v", c.ID, err)
		} else {
			zap.S().Infof("Sent disconnect message to client %s for room %s", c.ID, roomID)
		}
	} else {
		zap.S().Warnf("No available client found for room %s and client %s", roomID, c.ID)
	}
}

func (c *Client) NotifyClientNotInRoom(hub *Hub, roomID string) {
	msg := &ListenerMessage{
		Type: EventRoomUpdate,
		Data: &Message{
			ID:             uuid.NewString(),
			Text:           fmt.Sprintf("Client %s is not in room %s.", c.ID, roomID),
			RoomID:         roomID,
			CreatedAt:      time.Now(),
			ReceiverUserID: c.ID,
		},
		IsClient: false,
	}

	otherClient := hub.IsAvailableClient(roomID, c.ID)
	if otherClient != nil {
		if err := otherClient.Conn.WriteJSON(msg); err != nil {
			zap.S().Errorf("Error writing RoomNotIn message to client %s: %v", c.ID, err)
		} else {
			zap.S().Infof("Sent RoomNotIn message to client %s for room %s", c.ID, roomID)
		}
	} else {
		zap.S().Warnf("No available client found for room %s and client %s", roomID, c.ID)
	}
}
