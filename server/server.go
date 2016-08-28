package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Server represents a RESTful server
type Server struct {
	httpServer *http.Server

	// The router for HTTP service
	r *mux.Router
}

// New creates a new server
func New(conf *Config) (*Server, error) {
	if conf == nil {
		conf = &DefaultConfig
	}
	sb := serverBuilder{conf: *conf}
	server, err := sb.router().http().build()
	if err != nil {
		return nil, fmt.Errorf("set up server: %s", err)
	}
	return server, nil
}

// Start starts the server on the current process
func (s *Server) Start() {
	fmt.Println(s.httpServer.ListenAndServe())
}
