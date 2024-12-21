package chat

import (
	"errors"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gorilla/websocket"
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/http/middleware"
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/websocket/chat"
	"gitlab.otovinn.com/websocket-server/shared/data"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func (h *Handler) upgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := middleware.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		zap.S().Errorf(errConnectionUpgrade, err)
		http.Error(w, fmt.Sprintf(errConnectionUpgrade, err), http.StatusBadRequest)
		return nil, err
	}
	return conn, nil
}

func (h *Handler) isWSRoomAvailable(cl *chat.Client) error {
	if h.chatHub.IsCreatedRoom(cl.RoomID) {
		h.closeConnectionWithMessage(cl, websocket.ClosePolicyViolation, errMessageRead)
		return errors.New("room is already created")
	}
	return nil
}

func (h *Handler) createNewUser(req *CreateRoomReq) *chat.Client {
	return &chat.Client{
		Conn:    nil,
		Message: make(chan *chat.ListenerMessage, 100),
		RoomID:  req.UserID,
		ID:      req.UserID,
		Name:    req.Name,
		Status:  chat.Waiting,
		Role:    chat.User,
	}
}

func (h *Handler) createNewAdmin() *chat.Client {
	return &chat.Client{
		Message: make(chan *chat.ListenerMessage, 100),
		Role:    chat.Admin,
		Status:  chat.Waiting,
	}
}

func (h *Handler) closeConnectionWithMessage(cl *chat.Client, code int, message string) {
	closeMessage := websocket.FormatCloseMessage(code, message)
	cl.Conn.WriteMessage(websocket.CloseMessage, closeMessage)
	h.closeConnection(cl)
}

func (h *Handler) closeConnection(cl *chat.Client) {
	cl.Conn.Close()
	h.chatHub.Unregister <- cl
}

func (h *Handler) registerClient(cl *chat.Client) {
	h.chatHub.Register <- cl

	go cl.WriteMessage()
	cl.ReadMessage(h.chatHub)
}

func (h *Handler) authenticateClient(w http.ResponseWriter, r *http.Request, cl *chat.Client) bool {
	_, msg, err := cl.Conn.ReadMessage()
	if err != nil {
		h.closeConnectionWithMessage(cl, websocket.ClosePolicyViolation, errMessageRead)
		return false
	}

	var authMsg struct {
		Type  int    `json:"type"`
		Token string `json:"token"`
	}

	if err := json.Unmarshal(msg, &authMsg); err != nil {
		h.closeConnectionWithMessage(cl, websocket.ClosePolicyViolation, errInvalidAuthMessage)
		return false
	}

	if authMsg.Type != int(chat.EventAdminAuth) {
		h.closeConnectionWithMessage(cl, websocket.ClosePolicyViolation, errInvalidAuthMessageType)
		return false
	}

	user, err := h.middleware.ValidToken(r.Context(), authMsg.Token)
	if err != nil {
		h.closeConnectionWithMessage(cl, websocket.ClosePolicyViolation, errUnauthorizedToken)
		return false
	}

	cl.Name = user.FirstName + " " + user.LastName
	cl.ID = user.ID

	return true
}

func (h *Handler) checkShiftTime(w http.ResponseWriter, r *http.Request) error {
	now := time.Now().UTC().Add(3 * time.Hour)

	if !data.IsShiftTime(now) {
		http.Error(w, errShiftTimeInvalid, http.StatusBadRequest)
		return errors.New("invalid shift time")
	}

	return nil

}
