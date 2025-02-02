package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/cors"

	"github.com/gorilla/mux"

	"go.uber.org/zap"
)

type Server struct {
	ctx context.Context
	app *mux.Router
}

func NewServer(
	ctx context.Context,
) *Server {
	return &Server{
		ctx: ctx,
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

	serverConnURL := fmt.Sprintf("%s:%d", "localhost", 3000)

	srv := &http.Server{
		Addr:         serverConnURL,
		Handler:      c.Handler(s.app),
		ReadTimeout:  time.Minute * time.Duration(5),
		WriteTimeout: time.Minute * time.Duration(5),
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
