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

import "time"

const (
	// Done represents achievable done status
	Done = "DONE"
	// NotDone represents achievable not done status
	NotDone = "NOT_DONE"
	// InProgress represents achievable InProgress status
	InProgress = "IN_PROGRESS"
)

// Achievable represents an achievable action
// Achievable IS NOT a document in database, it is meant to be embedded
// inside other documents
type Achievable struct {
	Name        string    `bson:"name" json:"name" valid:"name" valid:"required,utfletternum,stringlength(1|100)"`
	Description string    `bson:"description,omitempty" json:"description,omitempty" valid:"optional,utfletternum,stringlength(1|400)"`
	LastUpdated time.Time `bson:"createdAt" json:"createdAt" valid:"-"`
	Status      string    `bson:"status" json:"status" valid:"validateStatus"`
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
