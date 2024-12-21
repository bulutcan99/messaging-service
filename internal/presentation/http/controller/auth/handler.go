package auth

import (
	"gitlab.otovinn.com/websocket-server/internal/core/application/usecase/auth"
	"gitlab.otovinn.com/websocket-server/internal/core/port"
)

type Handler struct {
	authService  port.AuthService
	tokenService *auth.TokenService
}

func NewHandler(authService port.AuthService, tokenService *auth.TokenService) *Handler {
	return &Handler{authService: authService, tokenService: tokenService}
}
