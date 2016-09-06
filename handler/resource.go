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

package handler

import (
	"context"
	"net/http"
	"path/filepath"
	"time"

	"github.com/iocat/donit/handler/internal/errors"
	"github.com/iocat/donit/handler/internal/service"
	"github.com/iocat/donit/handler/internal/utils"
	"github.com/iocat/donit/internal/donitdoc/followers"
	"github.com/iocat/donit/internal/donitdoc/users"

	docs "github.com/iocat/donit/internal/docgroup"
)

func init() {
	// TODO: set up services and database (service package)
}

// Endpoint serializes the HTTP endpoint
type Endpoint byte

const (
	// User represents the user endpoint
	User Endpoint = iota
	// Follower represents the follower endpoint
	Follower
	// Goal represents the goal endpoint
	Goal
	// Comment represents the comment endpoint
	Comment
	// Habit represents the habit endpoint
	Habit
	// Task represents the task endpoint
	Task
)

// URL gets the URL of the endpoint (resource)
func (e Endpoint) url() string {
	var endpointURL = []string{
		User:     "/users/{user}/goals/{goal}",
		Follower: "/users/{user}/followers/{follower}",
		Goal:     "/users/{user}/goals/{goal}",
		Comment:  "/users/{user}/goals/{goal}/comments/{comment}",
		Habit:    "/users/{user}/goals/{goal}/habits/{habit}",
		Task:     "/users/{user}/goals/{goal}/tasks/{task}",
	}
	return endpointURL[e]
}

// CollectionHandler ....
// TODO
func (e Endpoint) CollectionHandler() (string, http.HandlerFunc) {
	return e.baseURL(), UsersHandler
}

// ResourceHandler ....
// TODO
func (e Endpoint) ResourceHandler() (string, http.HandlerFunc) {
	return e.url(), UserHandler
}

// BaseURL gets the URL of the endpoint (collection)
func (e Endpoint) baseURL() string {
	return filepath.Dir(e.url())
}

var baseContext = context.Background()

// UsersHandler handles the operations on the user collection
func UsersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContextWithLog()
	switch r.Method {
	case "POST":
		var user docs.User
		if err := utils.DecodeJSON(r.Body, &user); err != nil {
			utils.HandleError(err, w)
			return
		}
		var getPassword = func(r *http.Request) (string, error) {
			if err := r.ParseForm(); err != nil {
				return "", errors.ErrBadData
			}
			password := r.Form.Get("password")
			if len(password) == 0 {
				return "", errors.NewBadData("password not provided")
			}
			return password, nil
		}
		password, err := getPassword(r)
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		err = service.User.Create(ctx, user.User, password)
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		utils.WriteJSONtoHTTP(nil, w, http.StatusCreated)
	default:
		utils.HandleError(errors.ErrMethodNotAllowed, w)
		return
	}
}

// UserHandler handles the operations on the user collection
func UserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContextWithLog()
	switch r.Method {
	case "GET":
		var user = docs.User{
			User: &users.User{},
		}
		// Check whether to expand the result
		toExpand, err := utils.ToExpand(r)
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		// Get the id
		ids, err := utils.MuxGetParams(r, "user")
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		username := ids[0]
		user.Username = username
		// Get the user
		if toExpand {
			err := user.Read(ctx)
			if err != nil {
				utils.HandleError(err, w)
				return
			}
		} else {
			user.User, err = service.User.Read(ctx, username)
			if err != nil {
				utils.HandleError(err, w)
				return
			}
		}
		utils.WriteJSONtoHTTP(user, w, http.StatusOK)
	case "PUT":
		var user = users.User{}
		if err := utils.DecodeJSON(r.Body, &user); err != nil {
			utils.HandleError(err, w)
			return
		}
		// Get the username
		ids, err := utils.MuxGetParams(r, "user")
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		user.Username = ids[0]
		// Update
		if err := service.User.Update(ctx, &user); err != nil {
			utils.HandleError(err, w)
			return
		}
		utils.WriteJSONtoHTTP(nil, w, http.StatusNoContent)
	case "DELETE":
		// Get the username
		ids, err := utils.MuxGetParams(r, "user")
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		err = service.User.Delete(ctx, ids[0])
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		utils.WriteJSONtoHTTP(nil, w, http.StatusNoContent)
	default:
		utils.HandleError(errors.ErrMethodNotAllowed, w)
		return
	}
}

// FollowersHandler handlers the operations on the followers collection
func FollowersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContextWithLog()
	// Get the username
	ids, err := utils.MuxGetParams(r, "user")
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	username := ids[0]
	switch r.Method {
	case "GET":
		lim, off, err := utils.GetLimitAndOffset(r)
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		fs, err := service.User.AllFollowers(ctx, username, lim, off)
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		utils.WriteJSONtoHTTP(fs, w, http.StatusOK)
	case "POST":
		follower := followers.Follower{}
		if err := utils.DecodeJSON(r.Body, follower); err != nil {
			utils.HandleError(err, w)
			return
		}
		follower.FollowAt, follower.Username = time.Now(), username
		err := service.Follower.Follows(ctx, &follower)
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		utils.WriteJSONtoHTTP(nil, w, http.StatusCreated)
	default:
		utils.HandleError(errors.ErrMethodNotAllowed, w)
		return
	}
}

// FollowerHandler handles the operations on the follower resource
func FollowerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := utils.NewContextWithLog()
	// Get the username
	ids, err := utils.MuxGetParams(r, "user", "follower")
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	username, follower := ids[0], ids[1]
	switch r.Method {
	case "DELETE":
		if err := service.Follower.Unfollows(ctx, username, follower); err != nil {
			utils.HandleError(err, w)
			return
		}
		utils.WriteJSONtoHTTP(nil, w, http.StatusNoContent)
	default:
		utils.HandleError(errors.ErrMethodNotAllowed, w)
		return
	}
}

// TODO:
// 		Write handlers for comments
//		Write handlers for goals
// 		Write handlers for habits
// 		Write handlers for tasks
