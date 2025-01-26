package server

import (
	"net/http"
)

func registerRoutes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users", handleGetUsers)
	mux.HandleFunc("POST /users", handleCreateUser)
	mux.HandleFunc("GET /feeds", handleGetFeeds)
	mux.HandleFunc("POST /feeds", handleCreateFeed)
	mux.HandleFunc("GET /hello-world", helloWorldHandler)

	// api
	mux.HandleFunc("GET /api/health", healthHandler)

	return mux
}
