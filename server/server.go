// Copyright 2016 Thanh Ngo <felix.infinite@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
