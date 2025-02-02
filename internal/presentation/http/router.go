package http

import (
	"net/http"
	"websocket-azure/presentation/http/controller/_default"
)

func (s *Server) SetupRouter() {
	s.index()
}

func (s *Server) index() {
	s.app.HandleFunc("/health", _default.HealthCheckHandler).Methods(http.MethodGet)
}

func (s *Server) message() {
	s.app.HandleFunc("/conversations/{conversationId}/messages", message.SendMessageHandler).Methods(http.MethodPost)
}
