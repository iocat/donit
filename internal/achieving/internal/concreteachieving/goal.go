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
func (cg *Goal) AddAchievable(achieving.Achievable) (utils.HexID, error) {
	return utils.HexID{}, nil
}

// RemoveAchievableTask removes the task
func (cg *Goal) RemoveAchievable(utils.HexID) error {
	return nil
}

// UpdateAchievableTask updates the task
func (cg *Goal) UpdateAchievable(achieving.Achievable, utils.HexID) error {
	return nil
}

// RetrieveAchievableTask retrieves the task list
func (cg *Goal) RetrieveAchievable() ([]achieving.Achievable, error) {
	return nil, nil
}
