package middleware

import (
	"gitlab.otovinn.com/websocket-server/internal/core/application/usecase/auth"
	"gitlab.otovinn.com/websocket-server/internal/core/port"
)

type Handler struct {
	tokenService *auth.TokenService
	authService  port.AuthService
}

func NewHandler(authService port.AuthService, tokenService *auth.TokenService) *Handler {
	return &Handler{tokenService: tokenService, authService: authService}
}
