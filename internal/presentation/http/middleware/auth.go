package middleware

import (
	"context"
	"errors"
	"gitlab.otovinn.com/websocket-server/internal/core/domain/entity"
	"net/http"
	"strings"
)

const (
	AuthHeader  = "Authorization"
	AuthType    = "Bearer"
	AuthSession = "Session"
)

func (h *Handler) ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get(AuthHeader)

		user, err := h.ValidToken(r.Context(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), AuthSession, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) ValidToken(ctx context.Context, token string) (entity.User, error) {
	if token == "" {
		return entity.User{}, errors.New("unauthorized")
	}

	fields := strings.Fields(token)
	if len(fields) != 2 {
		return entity.User{}, errors.New("unauthorized")
	}

	if fields[0] != AuthType {
		return entity.User{}, errors.New("unauthorized")
	}

	accessToken := fields[1]
	claims, err := h.tokenService.VerifyAccessToken(accessToken)
	if err != nil {
		return entity.User{}, errors.New("unauthorized")
	}

	user, err := h.authService.Me(ctx, claims.ID)
	if err != nil {
		return entity.User{}, errors.New("unauthorized")
	}

	return user, nil
}
