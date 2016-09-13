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
	"github.com/iocat/donit/internal/achieving/internal/user"
	"github.com/iocat/donit/internal/achieving/utils"
	"gopkg.in/mgo.v2"
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
func (c User) CreateGoal(g achieving.Goal) (utils.HexID, error) {
	if g, ok := g.(*Goal); ok {
		return c.User.CreateGoal(c.goalCollection, &(g.Goal))
	}
	return utils.HexID{}, fmt.Errorf("invalid data type, expect Goal, got %T", g)
}

// DeleteGoal deletes a goal
func (c User) DeleteGoal(id utils.HexID) error {
	return c.User.DeleteGoal(c.goalCollection, id)
}

// UpdateGoal updates a goal
func (c User) UpdateGoal(g achieving.Goal, id utils.HexID) error {
	if g, ok := g.(*Goal); ok {
		return c.User.UpdateGoal(c.goalCollection, &(g.Goal), id)
	}
	return fmt.Errorf("invalid data type, expect Goal, got %T", g)
}

// RetrieveGoal retrieves the goal
func (c User) RetrieveGoal(id utils.HexID) (achieving.Goal, error) {
	g, err := c.User.RetrieveGoal(c.goalCollection, id)
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
	gs, err := c.User.RetriveGoals(c.goalCollection, limit, offset)
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
