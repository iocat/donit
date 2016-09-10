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

// Package jsoninterpreter handles JSON marshaling and json unmarshaling
package jsoninterpreter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/iocat/donit/internal/achieving"
	concr "github.com/iocat/donit/internal/achieving/internal/concreteachieving"
	"gopkg.in/mgo.v2"
)

type errInvalidJSONType string

func (err errInvalidJSONType) Error() string {
	return string(err)
}

// IsErrInvalidJSONType checks whether the error is caused by
// wrong types
func IsErrInvalidJSONType(err error) bool {
	_, ok := err.(errInvalidJSONType)
	return ok
}

func newErrInvalidJSONType(err string) error {
	return errInvalidJSONType(err)
}

// Interpreter represents a JSON interpreter which  decodes and encodes object
type Interpreter interface {
	// Decode interprets and decodes the JSON data into the output object
	Decode(io.Reader) (interface{}, error)

	// Encode writes the object to the io.Writer output buffer
	//
	// If the object implements json.Marshaler, the marshalling output will be
	// written instead
	Encode(io.Writer, interface{}) error
}

// UserJSONInterpreter implements Interpreter
type UserJSONInterpreter struct {
	goal       *mgo.Collection
	achievable *mgo.Collection
}

// NewUser creates a user interpreter
func NewUser(user, goal, achievable *mgo.Collection) Interpreter {
	return UserJSONInterpreter{
		goal:       goal,
		achievable: achievable,
	}
}

func (u UserJSONInterpreter) decode(r io.Reader) (achieving.User, error) {
	var user = concr.NewUser(u.goal, u.achievable)
	if err := json.NewDecoder(r).Decode(&user.User); err != nil {
		return nil, newErrInvalidJSONType(err.Error())
	}
	return user, nil
}

// Decode implements Interpreter's Decode
func (u UserJSONInterpreter) Decode(r io.Reader) (interface{}, error) {
	return u.decode(r)
}

func (u UserJSONInterpreter) encode(w io.Writer, user achieving.User) error {
	if err := json.NewEncoder(w).Encode(user); err != nil {
		return err
	}
	return nil
}

// Encode implements Interpreter's Encode
func (u UserJSONInterpreter) Encode(w io.Writer, user interface{}) error {
	casted, ok := user.(achieving.User)
	if !ok {
		return newErrInvalidJSONType(fmt.Sprintf("the provided type is not User, got %T", user))
	}
	return u.encode(w, casted)
}

// GoalJSONInterpreter interprets goal Json request
type GoalJSONInterpreter struct {
	achievable *mgo.Collection
}

// NewGoal creates a new interpreter for goal
func NewGoal(achievable *mgo.Collection) Interpreter {
	return GoalJSONInterpreter{
		achievable: achievable,
	}
}

func (g GoalJSONInterpreter) decode(r io.Reader) (achieving.Goal, error) {
	var goal = concr.NewGoal(g.achievable)
	if err := json.NewDecoder(r).Decode(&(goal.Goal)); err != nil {
		return nil, newErrInvalidJSONType(err.Error())
	}
	return goal, nil
}

// Decode implements Interpreter's Decode
func (g GoalJSONInterpreter) Decode(r io.Reader) (interface{}, error) {
	return g.decode(r)
}

func (g GoalJSONInterpreter) encode(w io.Writer, goal achieving.Goal) error {
	if err := json.NewEncoder(w).Encode(goal); err != nil {
		return err
	}
	return nil
}

// Encode implements Interpreter's Encode
func (g GoalJSONInterpreter) Encode(w io.Writer, goal interface{}) error {
	casted, ok := goal.(achieving.Goal)
	if !ok {
		return newErrInvalidJSONType(fmt.Sprintf("the provided type is not Goal, got %T", goal))
	}
	return g.encode(w, casted)
}

// NewAchievable creates a new interpreter for Achievable task
func NewAchievable() Interpreter {
	return AchievableJSONInterpreter{}
}

// AchievableJSONInterpreter represents a JSON decoder/encoder for Achievable tasks
type AchievableJSONInterpreter struct{}

func (a AchievableJSONInterpreter) decode(r io.Reader) (achieving.Achievable, error) {
	var ach = concr.Achievable{}
	if err := json.NewDecoder(r).Decode(&(ach.Achievable)); err != nil {
		return nil, newErrInvalidJSONType(err.Error())
	}
	return &ach, nil
}

// Decode implements Interpreter's Decode
func (a AchievableJSONInterpreter) Decode(r io.Reader) (interface{}, error) {
	return a.decode(r)
}

func (a AchievableJSONInterpreter) encode(w io.Writer, ach achieving.Achievable) error {
	if err := json.NewEncoder(w).Encode(ach); err != nil {
		return err
	}
	return nil
}

// Encode implements Interpreter's Encode
func (a AchievableJSONInterpreter) Encode(w io.Writer, ach interface{}) error {
	casted, ok := ach.(achieving.Achievable)
	if !ok {
		return newErrInvalidJSONType(fmt.Sprintf("the provided type is not Achievable, got %T", ach))
	}
	return a.encode(w, casted)
}
