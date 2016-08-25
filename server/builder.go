package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/iocat/donit/handler"
)

type serverBuilder struct {
	Server
	conf Config
}

func (sb *serverBuilder) build() *Server {
	return &sb.Server
}

// setupRouter sets up the router with a prefedefined path
func (sb *serverBuilder) router() *serverBuilder {
	sb.r = mux.NewRouter()
	// Set up a subrouter for api that matches to "http[s]://api.domain.com/v#/"
	common := sb.r
	common.NotFoundHandler = handler.NotFound
	// Asign handlers

	common.HandleFunc("/users", handler.Get("users"))
	common.HandleFunc("/users/{user}", handler.Get("user"))
	common.HandleFunc("/users/{user}/validate", handler.Get("validator"))

	common.HandleFunc("/users/{user}/goals", handler.Get("goals"))
	common.HandleFunc("/users/{user}/goals/{goal}", handler.Get("goal"))

	common.HandleFunc("/users/{user}/goals/{goal}/habits", handler.Get("habits"))
	common.HandleFunc("/users/{user}/goals/{goal}/habits/{habit}", handler.Get("habit"))

	common.HandleFunc("/users/{user}/goals/{goal}/tasks", handler.Get("tasks"))
	common.HandleFunc("/users/{user}/goals/{goal}/tasks/{task}", handler.Get("task"))
	return sb
}

// http sets up the http server
// TODO: add TLS layer for production usage
func (sb *serverBuilder) http() *serverBuilder {
	sb.httpServer = &http.Server{
		Handler:      sb.r,
		Addr:         "127.0.0.1:8080",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	sb.httpServer.SetKeepAlivesEnabled(false)
	return sb
}
