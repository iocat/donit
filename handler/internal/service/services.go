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

package service

import (
	"fmt"

	"gopkg.in/mgo.v2"

	"github.com/go-kit/kit/log"

	"github.com/iocat/donit/internal/achiever"
	"github.com/iocat/donit/internal/achieving"
)

var (
	User     achiever.UserService
	Follower achiever.FollowerService
	Goal     achieving.GoalService
	Task     achieving.TaskService
	Habit    achieving.HabitService
	Comment  achieving.CommentService
)

// SetUp sets up all the data service
func SetUp(l log.Logger, db *mgo.Database) error {
	var err error
	User, err = achiever.NewUser(l, db)
	if err != nil {
		return fmt.Errorf("setting up service: %s", err)
	}
	Follower, err = achiever.NewFollower(l, db)
	if err != nil {
		return fmt.Errorf("setting up service: %s", err)
	}
	Goal, err = achieving.NewGoal(l, db)
	if err != nil {
		return fmt.Errorf("setting up service: %s", err)
	}
	Task, err = achieving.NewTask(l, db)
	if err != nil {
		return fmt.Errorf("setting up service: %s", err)
	}
	Habit, err = achieving.NewHabit(l, db)
	if err != nil {
		return fmt.Errorf("setting up service: %s", err)
	}
	Comment, err = achieving.NewComment(l, db)
	if err != nil {
		return fmt.Errorf("setting up service: %s", err)
	}
	return nil
}
