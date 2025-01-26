package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", s.handleGetUsers)
	mux.HandleFunc("POST /users", s.handleCreateUser)
	mux.HandleFunc("GET /feeds", s.handleGetFeeds)
	mux.HandleFunc("POST /feeds", s.handleCreateFeed)
	mux.HandleFunc("GET /hello-world", s.helloWorldHandler)

	// api
	mux.HandleFunc("GET /api/health", s.healthHandler)

	return mux
}
