package chat

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
)

type CreateRoomReq struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
	Plate  string `json:"plate"`
}

func (h *Handler) UserJoinRoom(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	userID := mux.Vars(r)["user_id"]
	plate := mux.Vars(r)["plate"]

	req := &CreateRoomReq{
		Name:   name,
		UserID: userID,
		Plate:  plate,
	}

	cl := h.createNewUser(req)

	if err := h.checkShiftTime(w, r); err != nil {
		zap.S().Errorf("Shift time check failed: %v", err)
		http.Error(w, "Shift time error", http.StatusForbidden)
		return
	}

	err := h.CreateRoomHandler(cl)
	if err != nil {
		zap.S().Errorf("Failed to create room handler: %v", err)
		http.Error(w, "Create room error", http.StatusForbidden)
		return
	}

	conn, err := h.upgradeConnection(w, r)
	if err != nil {
		zap.S().Errorf("Failed to upgrade connection: %v", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	cl.Conn = conn

	h.registerClient(cl)
}
