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

package concreteachieving

import (
	"fmt"

	"github.com/iocat/donit/internal/achieving"
	"github.com/iocat/donit/internal/achieving/errors"
	"github.com/iocat/donit/internal/achieving/internal/user"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// NewUser creates a new empty User
func NewUser(goal, task *mgo.Collection) *User {
	return &User{
		goalCollection:       goal,
		achievableCollection: task,
	}
}

// User represents the concrete user
type User struct {
	user.User            `valid:"required"`
	achievableCollection *mgo.Collection `valid:"-"`
	goalCollection       *mgo.Collection `valid:"-"`
}

// CreateGoal creates a new goal
func (c User) CreateGoal(g achieving.Goal) (string, error) {
	if g, ok := g.(*Goal); ok {
		id, err := c.User.CreateGoal(c.goalCollection, &(g.Goal))
		if err != nil {
			return "", err
		}
		return id.Hex(), nil
	}
	return "", fmt.Errorf("invalid data type, expect Goal, got %T", g)
}

// DeleteGoal deletes a goal
func (c User) DeleteGoal(id string) error {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return errors.NewValidate(fmt.Sprintf("%s is not a valid resource id", id))
	}
	return c.User.DeleteGoal(c.goalCollection, bson.ObjectIdHex(id))
}

// UpdateGoal updates a goal
func (c User) UpdateGoal(g achieving.Goal, id string) error {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return errors.NewValidate(fmt.Sprintf("%s is not a valid resource id", id))
	}
	if g, ok := g.(*Goal); ok {
		return c.User.UpdateGoal(c.goalCollection, &(g.Goal), bson.ObjectIdHex(id))
	}
	return fmt.Errorf("invalid data type, expect Goal, got %T", g)
}

// RetrieveGoal retrieves the goal
func (c User) RetrieveGoal(id string) (achieving.Goal, error) {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return nil, errors.NewValidate(fmt.Sprintf("%s is not a valid resource id", id))
	}
	g, err := c.User.RetrieveGoal(c.goalCollection, c.achievableCollection, bson.ObjectIdHex(id))
	if err != nil {
		return nil, err
	}
	return &Goal{
		Goal:                 g,
		achievableCollection: c.achievableCollection,
	}, nil
}

// RetrieveGoals retrieves a goal
func (c User) RetrieveGoals(limit, offset int) ([]achieving.Goal, error) {
	gs, err := c.User.RetriveGoals(c.goalCollection, c.achievableCollection, limit, offset)
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
