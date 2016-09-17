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
	"time"

	"github.com/gorilla/mux"
	"github.com/iocat/donit/handler"
)

type serverBuilder struct {
	Server
	conf Config
	err  error
}

func (sb *serverBuilder) build() (*Server, error) {
	if sb.err != nil {
		return nil, sb.err
	}
	return &sb.Server, nil
}

// setupRouter sets up the router with a prefedefined path
func (sb *serverBuilder) router() *serverBuilder {
	sb.r = mux.NewRouter()
	// Set up a handler dispatcher
	common := sb.r
	common.NotFoundHandler = handler.NotFound
	common.HandleFunc(handler.User.BaseURL(), handler.CreateUser).Methods("POST")
	common.HandleFunc(handler.User.URL(), handler.DeleteUser).Methods("DELETE")
	common.HandleFunc(handler.User.URL(), handler.UpdateUser).Methods("PUT")
	common.HandleFunc(handler.User.URL(), handler.ReadUser).Methods("GET")

	common.HandleFunc(handler.Goal.BaseURL(), handler.CreateGoal).Methods("POST")
	common.HandleFunc(handler.Goal.BaseURL(), handler.AllGoals).Methods("GET")
	common.HandleFunc(handler.Goal.URL(), handler.DeleteGoal).Methods("DELETE")
	common.HandleFunc(handler.Goal.URL(), handler.UpdateGoal).Methods("PUT")
	common.HandleFunc(handler.Goal.URL(), handler.ReadGoal).Methods("GET")

	common.HandleFunc(handler.Achievable.BaseURL(), handler.CreateAchievable).Methods("POST")
	common.HandleFunc(handler.Achievable.BaseURL(), handler.AllAchievables).Methods("GET")
	common.HandleFunc(handler.Achievable.URL(), handler.DeleteAchievable).Methods("DELETE")
	common.HandleFunc(handler.Achievable.URL(), handler.UpdateAchievable).Methods("PUT")

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
