package chat

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
)

const (
	errRoomNotCreated         = "Room %s is not created"
	errConnectionUpgrade      = "Error upgrading connection: %v"
	errMessageRead            = "Error reading message"
	errInvalidAuthMessage     = "Invalid auth message format"
	errInvalidAuthMessageType = "Invalid auth message type"
	errUnauthorizedToken      = "Unauthorized: Invalid token"
	errShiftTimeInvalid       = "Shift time is not valid"
	errJSONMarshal            = "JSON marshal error"
)

func (h *Handler) AdminJoinHub(w http.ResponseWriter, r *http.Request) {
	if err := h.checkShiftTime(w, r); err != nil {
		zap.S().Errorf("Shift time check failed: %v", err)
		http.Error(w, "Shift time error", http.StatusForbidden)
		return
	}

	conn, err := h.upgradeConnection(w, r)
	if err != nil {
		zap.S().Errorf("Failed to upgrade connection: %v", err)
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}

	cl := h.createNewAdmin()
	cl.Conn = conn

	if !h.authenticateClient(w, r, cl) {
		zap.S().Error("Failed to authenticate admin client")
		h.closeConnectionWithMessage(cl, websocket.ClosePolicyViolation, errInvalidAuthMessage)
		return
	}

	h.registerClient(cl)
}
