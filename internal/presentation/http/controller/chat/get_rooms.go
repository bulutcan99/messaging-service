package chat

import (
	"github.com/gorilla/mux"
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/http/dto"
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/http/response"
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/websocket/chat"
	"net/http"
)

func (h *Handler) GetRoom(w http.ResponseWriter, r *http.Request) {
	roomID := mux.Vars(r)["room_id"]

	if roomInterface, ok := h.chatHub.Rooms.Load(roomID); ok {
		room := roomInterface.(*chat.Room)
		clients := h.GetClientsForRoom(room.ID)

		roomRes := dto.RoomRes{
			ID:    room.ID,
			Name:  room.Name,
			Users: clients,
		}

		response.RespondWithJSON(w, http.StatusOK, map[string]any{
			"error": false,
			"room":  roomRes,
		})
		return
	}

	response.RespondWithJSON(w, http.StatusNotFound, map[string]any{
		"error":   true,
		"message": "Room not found",
	})
}

func (h *Handler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms := make([]dto.RoomRes, 0)

	h.chatHub.Rooms.Range(func(key, value interface{}) bool {
		room := value.(*chat.Room)
		clients := h.GetClientsForRoom(room.ID)

		rooms = append(rooms, dto.RoomRes{
			ID:    room.ID,
			Name:  room.Name,
			Users: clients,
		})
		return true
	})

	response.RespondWithJSON(w, http.StatusOK, map[string]any{
		"error": false,
		"rooms": rooms,
	})
}

func (h *Handler) GetClientsForRoom(roomID string) []dto.ClientRes {
	var clients []dto.ClientRes

	if roomInterface, ok := h.chatHub.Rooms.Load(roomID); ok {
		room := roomInterface.(*chat.Room)

		for _, c := range room.Clients {
			clients = append(clients, dto.ClientRes{
				ID:   c.ID,
				Name: c.Name,
			})
		}
	}

	return clients
}
