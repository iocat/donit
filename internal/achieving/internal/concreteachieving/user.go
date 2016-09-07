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

// Package concreteachieving contains concrete implementation of the achieving package
// concreteachievable encapsulates the creation of the private concrete
// implemenation, so that the jsoninterpreter can creates object. (All creation
// of object should be via the jsoninterpreter)
package concreteachieving

import (
	"encoding/json"

	"github.com/go-errors/errors"
	"github.com/iocat/donit/internal/achieving"
	"github.com/iocat/donit/internal/achieving/internal/user"
	"github.com/iocat/donit/internal/achieving/utils"
	"gopkg.in/mgo.v2"
)

// NewUser creates a new empty User
func NewUser(goalCollection *mgo.Collection, achievableCollection *mgo.Collection) achieving.User {
	return &User{
		goalCollection:       goalCollection,
		achievableCollection: achievableCollection,
	}
}

// NewUserFromJSON possibly returns a json decoding error
func NewUserFromJSON(d *json.Decoder, goalCollection *mgo.Collection, achievableCollection *mgo.Collection) (achieving.User, error) {
	var u User
	err := d.Decode(&(u.User))
	if err != nil {
		return nil, errors.Errorf("invalid json")
	}
	u.goalCollection, u.achievableCollection = goalCollection, achievableCollection
	return u, nil
}

// User represents the concrete user
type User struct {
	user.User
	achievableCollection *mgo.Collection
	goalCollection       *mgo.Collection
}

// CreateGoal creates a new goal
func (c User) CreateGoal(g achieving.Goal) (utils.HexID, error) {
	return c.User.CreateGoal(c.goalCollection, &(g.(*Goal)).Goal)
}

// DeleteGoal deletes a goal
func (c User) DeleteGoal(id utils.HexID) error {
	return c.User.DeleteGoal(c.goalCollection, id)
}

// UpdateGoal updates a goal
func (c User) UpdateGoal(g achieving.Goal, id utils.HexID) error {
	return c.User.UpdateGoal(c.goalCollection, &(g.(*Goal)).Goal, id)
}

// RetrieveGoals retrieves a goal
func (c User) RetrieveGoals() ([]achieving.Goal, error) {
	gs, err := c.User.RetriveGoals(c.goalCollection)
	if err != nil {
		return nil, err
	}
	var goals []achieving.Goal
	for _, g := range gs {
		goals = append(goals, &Goal{
			Goal:                 g,
			achievableCollection: c.achievableCollection,
		})
	}
	return goals, nil
}
