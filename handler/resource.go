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
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/iocat/donit/internal/achieving"
	json "github.com/iocat/donit/internal/achieving/jsoninterpreter"
)

// Endpoint serializes the HTTP endpoint
type Endpoint byte

const (
	// User represents the user endpoint
	User Endpoint = iota
	// Goal represents the goal endpoint
	Goal
	// Achievable represents an achievable task endpoint
	Achievable
)

// URL gets the URL of the endpoint (resource)
func (e Endpoint) URL() string {
	return e.url()
}

func (e Endpoint) url() string {
	var endpointURL = []string{
		User:       "/users/{user}",
		Goal:       "/users/{user}/goals/{goal}",
		Achievable: "/users/{user}/goals/{goal}/achievables/{achievable}",
	}
	return endpointURL[e]
}

// BaseURL gets the URL of the endpoint (collection)
func (e Endpoint) BaseURL() string {
	return e.baseURL()
}

func (e Endpoint) baseURL() string {
	return filepath.Dir(e.URL())
}

func (e Endpoint) resourceKeyNames() []string {
	var keyNames = [][]string{
		User:       []string{"user"},
		Goal:       []string{"user", "goal"},
		Achievable: []string{"user", "achievable"},
	}
	return keyNames[e]
}

func (e Endpoint) collectionKeyNames() []string {
	rkn := e.resourceKeyNames()
	return rkn[:len(rkn)-1]
}

// TODO: not implemented
func (e Endpoint) resourceLocationForID(id ...string) string {

	return ""
}

var baseContext = context.Background()

func getResourceFromContext(ctx context.Context) interface{} {
	return ctx.Value("resource")
}

var store achieving.UserStore

var collections = []*mgo.Collection{
	User:       nil,
	Goal:       nil,
	Achievable: nil,
}

var interpreters = []json.Interpreter{
	User:       nil,
	Goal:       nil,
	Achievable: nil,
}

func (e Endpoint) collection() *mgo.Collection {
	return collections[e]
}

func (e Endpoint) interpreter() json.Interpreter {
	return interpreters[e]
}

var userInterpreter, goalInterpreter, achievableInterpreter json.Interpreter

func init() {
	sess, err := mgo.DialWithTimeout("localhost:27017", 5*time.Second)
	if err != nil {
		fmt.Printf("set up database: %s", err)
		os.Exit(1)
	}
	db := sess.DB("donit")
	collections = []*mgo.Collection{
		User:       db.C("users"),
		Goal:       db.C("goals"),
		Achievable: db.C("achievables"),
	}
	interpreters = []json.Interpreter{
		User:       json.NewUser(Goal.collection(), Achievable.collection()),
		Goal:       json.NewGoal(Achievable.collection()),
		Achievable: json.NewAchievable(),
	}
	store = json.NewStore(User.collection(), Goal.collection(), Achievable.collection())
}
