package message

import "websocket-azure/shared/auth"

type Handler struct {
	tokenService *auth.TokenService
}

func NewHandler(tokenService *auth.TokenService) *Handler {
	return &Handler{
		tokenService,
	}
}
