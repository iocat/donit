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

package achievable

import (
	"time"

	"github.com/iocat/donit/internal/achieving/utils"
	"gopkg.in/mgo.v2/bson"
)

// Reminder represents a reminder for tasks
type Reminder struct {
	At       time.Time     `bson:"remindAt" json:"remindAt" valid:"-"`
	Duration time.Duration `bson:"duration" json:"duration" valid:"-"`
}

// Task represents a task
type Task struct {
	utils.HexID `bson:"_id,omitempty" json:"id" valid:"required,hexadecimal"`
	Goal        bson.ObjectId `bson:"goal" json:"-" valid:"required, hexadecimal"`
	Achievable  `bson:"subGoal,inline" valid:"required"`
	*Reminder   `bson:"reminder,omitempty" json:"reminder,omitempty" valid:"optional"`
}
