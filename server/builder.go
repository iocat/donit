package server

import (
	"fmt"
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
	var resources = handler.GetResources()
	for _, h := range resources {
		col, ite := h.URL()
		common.HandleFunc(col, h.Collection())
		common.HandleFunc(ite, h.Item())
	}
	common.HandleFunc("/users/{user}/validate", handler.Get("validator"))
	return sb
}

// http sets up the http server
// TODO: add TLS layer for production usage
func (sb *serverBuilder) http() *serverBuilder {
	sb.httpServer = &http.Server{
		Handler:      sb.r,
		Addr:         fmt.Sprintf("%s:%d", sb.conf.Domain, sb.conf.Port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	sb.httpServer.SetKeepAlivesEnabled(false)
	return sb
}
