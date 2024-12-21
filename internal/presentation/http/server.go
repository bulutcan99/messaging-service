package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/cors"
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/http/controller/auth"
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/http/controller/chat"
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/http/middleware"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"gitlab.otovinn.com/websocket-server/shared/config"
	"go.uber.org/zap"
)

type Server struct {
	ctx               context.Context
	cfg               *config.Container
	app               *mux.Router
	middlewareHandler *middleware.Handler
	authHandler       *auth.Handler
	chatHandler       *chat.Handler
}

func NewServer(
	ctx context.Context,
	cfg *config.Container,
	middlewareHandler *middleware.Handler,
	authHandler *auth.Handler,
	chatHandler *chat.Handler,

) *Server {
	return &Server{
		ctx:               ctx,
		cfg:               cfg,
		middlewareHandler: middlewareHandler,
		authHandler:       authHandler,
		chatHandler:       chatHandler,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.app = mux.NewRouter()
	s.SetupRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	serverConnURL := fmt.Sprintf("%s:%d", s.cfg.HTTP.Host, s.cfg.HTTP.Port)

	srv := &http.Server{
		Addr:         serverConnURL,
		Handler:      c.Handler(s.app),
		ReadTimeout:  time.Minute * time.Duration(s.cfg.Settings.ServerReadTimeout),
		WriteTimeout: time.Minute * time.Duration(s.cfg.Settings.ServerReadTimeout),
		IdleTimeout:  time.Minute * 2,
	}

	go func() {
		zap.S().Debug("Starting HTTP server on ", serverConnURL)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.S().Fatalf("server listen error: %v", err)
		}
	}()

	return nil
}

func (s *Server) Close(ctx context.Context) error {
	zap.S().Debug("HTTP-Server Context is done. Shutting down server...")
	srv := &http.Server{
		Handler: s.app,
	}

	if err := srv.Shutdown(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			return nil
		}
		zap.S().Errorf("server shutdown error: %v", err)
		return err
	}

	return nil
}
