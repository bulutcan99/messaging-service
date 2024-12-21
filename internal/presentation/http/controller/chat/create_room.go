package chat

import (
	"errors"
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/websocket/chat"
)

func (h *Handler) CreateRoomHandler(cl *chat.Client) error {
	if h.chatHub.IsCreatedRoom(cl.RoomID) {
		return errors.New("room is already created")
	}

	h.chatHub.RoomRegister <- cl
	return nil
}
