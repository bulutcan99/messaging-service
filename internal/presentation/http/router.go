package http

import (
	"gitlab.otovinn.com/websocket-server/internal/adapter/driving/presentation/http/controller/_default"
	"net/http"
)

func (s *Server) SetupRouter() {
	s.index()
	s.admin()
	s.user()
}

func (s *Server) index() {
	s.app.HandleFunc("/health", _default.HealthCheckHandler).Methods(http.MethodGet)
}

func (s *Server) admin() {
	s.app.HandleFunc("/api/v1/admin/login", s.authHandler.AdminLoginHandler).Methods(http.MethodPost)
	s.app.Handle("/api/v1/admin/hub/join", http.HandlerFunc(s.chatHandler.AdminJoinHub))
	//s.app.Handle("/api/v1/admin/rooms/{roomID}/messages", s.middlewareHandler.ValidateJWT(http.HandlerFunc(s.chatHandler.GetRoomMessages))).Methods("GET")
}

func (s *Server) user() {
	s.app.Handle("/api/v1/user/{name}:{user_id}:{plate}/rooms/join", http.HandlerFunc(s.chatHandler.UserJoinRoom))
}
