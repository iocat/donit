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
	"github.com/iocat/donit/internal/achieving/internal/goal"
	"github.com/iocat/donit/internal/achieving/utils"
	"gopkg.in/mgo.v2"
)

// Goal implements the achieving.Goal interface
type Goal struct {
	goal.Goal
	achievableCollection *mgo.Collection
}

// AddAchievable adds a new achievable task
func (cg *Goal) AddAchievable(a achieving.Achievable) (utils.HexID, error) {
	if a, ok := a.(*Achievable); ok {
		return cg.Goal.AddAchievable(cg.achievableCollection, &(a.Achievable))
	}
	return utils.HexID{}, fmt.Errorf("wrong data type, expect Achievable, got %T", a)
}

// RemoveAchievable removes the task
func (cg *Goal) RemoveAchievable(id utils.HexID) error {
	return cg.RemoveAchievable(id)
}

// UpdateAchievable updates the task
func (cg *Goal) UpdateAchievable(a achieving.Achievable, id utils.HexID) error {
	if a, ok := a.(*Achievable); ok {
		return cg.Goal.UpdateAchievable(cg.achievableCollection, &(a.Achievable), id)
	}
	return fmt.Errorf("wrong data type, expect Achievable, got %T", a)
}

// RetrieveAchievable retrieves the task list
func (cg *Goal) RetrieveAchievable(limit, offset int) ([]achieving.Achievable, error) {
	as, err := cg.Goal.RetrieveAchievable(cg.achievableCollection, limit, offset)
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
