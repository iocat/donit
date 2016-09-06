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

package goals

import (
	"github.com/iocat/donit/internal/donitdoc/achievable"
	valid "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

const (
	// AccessPrivate is private accessibility
	AccessPrivate = "PRIVATE"
	// AccessForFollowers is the accessibility for followers
	AccessForFollowers = "FOR_FOLLOWERS"
	// AccessPublic is the accessibility for public user
	AccessPublic = "PUBLIC"
)

func init() {
	valid.SetFieldsRequiredByDefault(true)
	valid.CustomTypeTagMap.Set("goalAccessValidator", valid.CustomTypeValidator(GoalAccessValidatorFunc))
	valid.CustomTypeTagMap.Set("validateStatus", valid.CustomTypeValidator(achievable.ValidateStatus))
}

// Goal represents an achievable Goal
type Goal struct {
	bson.ObjectId         `bson:"_id,omitempty" json:"id" valid:"required,hexadecimal"`
	Username              string `bson:"username" json:"username" valid:"required,alphanum,length(1|30)"`
	achievable.Achievable `bson:"subGoal,inline" valid:"required"`
	PictureURL            string `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty" valid:"optional,url"`
	Accessibility         string `bson:"accessibility" json:"accessibility,omitempty" valid:"required,goalAccessValidator"`
}

// GoalAccessValidatorFunc validates the accessibility field of the Goal model
func GoalAccessValidatorFunc(value, _ interface{}) bool {
	switch value := value.(type) {
	case string:
		switch value {
		case AccessPrivate, AccessPublic, AccessForFollowers:
			return true
		default:
			return false
		}
	default:
		panic("the accessibility field must be a string")
	}
}
