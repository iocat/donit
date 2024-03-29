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
	"github.com/iocat/donit/internal/achieving/internal/goal"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Goal implements the achieving.Goal interface
type Goal struct {
	goal.Goal            `valid:"required"`
	achievableCollection *mgo.Collection `valid:"-"`
}

// NewGoal creates a new goal
func NewGoal(taskCol *mgo.Collection) *Goal {
	return &Goal{
		achievableCollection: taskCol,
	}
}

// AddAchievable adds a new achievable task
func (cg *Goal) AddAchievable(a achieving.Achievable) (string, error) {
	if a, ok := a.(*Achievable); ok {
		id, err := cg.Goal.AddAchievable(cg.achievableCollection, &(a.Achievable))
		if err != nil {
			return "", err
		}
		return id.Hex(), nil
	}
	return "", fmt.Errorf("wrong data type, expect Achievable, got %T", a)
}

// RemoveAchievable removes the task
func (cg *Goal) RemoveAchievable(id string) error {
	ok := bson.IsObjectIdHex(id)
	if !ok {
		return errors.NewValidate(fmt.Sprintf("%s is not a valid resource id", id))
	}
	return cg.Goal.RemoveAchievable(cg.achievableCollection, bson.ObjectIdHex(id))
}

// UpdateAchievable updates the task
func (cg *Goal) UpdateAchievable(a achieving.Achievable, id string) error {
	if a, ok := a.(*Achievable); ok {
		ok := bson.IsObjectIdHex(id)
		if !ok {
			return errors.NewValidate(fmt.Sprintf("%s is not a valid resource id", id))
		}
		return cg.Goal.UpdateAchievable(cg.achievableCollection, &(a.Achievable),
			bson.ObjectIdHex(id))
	}
	return fmt.Errorf("wrong data type, expect Achievable, got %T", a)
}

// RetrieveAchievables retrieves the task list
func (cg *Goal) RetrieveAchievables(limit, offset int) ([]achieving.Achievable, error) {
	as, err := cg.Goal.RetrieveAchievables(cg.achievableCollection, limit, offset)
	if err != nil {
		return nil, err
	}
	var res []achieving.Achievable
	for _, a := range as {
		res = append(res, &Achievable{
			Achievable: a,
		})
	}
	return res, nil
}
