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

// Package achievable contains goal's achievable task data
package achievable

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	// Done represents achievable done status
	Done = "DONE"
	// NotDone represents achievable not done status
	NotDone = "NOT_DONE"
	// InProgress represents achievable InProgress status
	InProgress = "IN_PROGRESS"
)

// Reminder represents a reminder for tasks
type Reminder struct {
	At       time.Time     `bson:"remindAt" json:"remindAt" valid:"-"`
	Duration time.Duration `bson:"duration" json:"duration" valid:"-"`
}

// Achievable represents an achievable action
type Achievable struct {
	ID             bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty" valid:"optional,hexadecimal"`
	Goal           bson.ObjectId   `bson:"_goal,omitempty" json:"-" valid:"optional"`
	Name           string          `bson:"name" json:"name" valid:"name" valid:"required,utfletternum,stringlength(1|100)"`
	Description    string          `bson:"description,omitempty" json:"description,omitempty" valid:"optional,stringlength(1|400)"`
	Status         string          `bson:"status" json:"status" valid:"validateStatus"`
	Reminder       *Reminder       `bson:"reminder,omitempty" json:"reminder,omitempty" valid:"optional"`
	RepeatReminder *RepeatReminder `bson:"repreatedReminder,omitempty" json:"repeatedReminder,omitempty" valid:"optional"`
}

// IsHabit returns whether this is a habit
func (a *Achievable) IsHabit() bool {
	return a.RepeatReminder != nil
}

// IsTask returns whether this is a task or not
func (a *Achievable) IsTask() bool {
	return a.RepeatReminder == nil
}

// ValidateStatus validates the status field
func ValidateStatus(value, _ interface{}) bool {
	switch value := value.(type) {
	case string:
		switch value {
		case Done, NotDone, InProgress:
			return true
		default:
			return false
		}
	default:
		panic("the status field must be a string")
	}
}

// HasAchieved returns whether achievable object is achieved or not
func (a Achievable) HasAchieved() bool {
	return a.Status == Done
}
