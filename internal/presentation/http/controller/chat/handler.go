package chat

import (
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/http/middleware"
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/websocket/chat"
	"gitlab.otovinn.com/websocket-server/internal/core/port"
)

type Handler struct {
	middleware  *middleware.Handler
	chatHub     *chat.Hub
	authService port.AuthService
}

func NewHandler(authService port.AuthService, middleware *middleware.Handler) *Handler {
	hub := chat.NewHub()
	go hub.Run()

	return &Handler{
		chatHub:     hub,
		authService: authService,
		middleware:  middleware,
	}
}
