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
	"fmt"
	"time"
)

const (
	// EveryDay is the Cycle for Every day habit
	EveryDay = "EVERY_DAY"
	// EveryWeekAndCustom is the cycle for every week habit
	EveryWeekAndCustom = "EVERY_WEEK"
	// EveryMonthAndCustom is the cycle for every month habit
	EveryMonthAndCustom = "EVERY_MONTH"
)
const (
	// Sunday represents Sunday
	Sunday = iota
	// Monday represents Monday
	Monday
	// Tuesday represents Tuesday
	Tuesday
	// Wednesday represents Wednesday
	Wednesday
	// Thursday represents Thursday
	Thursday
	// Friday represents Friday
	Friday
	// Saturday represents Saturday
	Saturday
)

// ValidateCycle validates the Cycle field
func ValidateCycle(v, _ interface{}) bool {
	switch v := v.(type) {
	case string:
		switch v {
		case EveryDay, EveryMonthAndCustom, EveryWeekAndCustom:
			return true
		default:
			return false
		}
	default:
		panic(fmt.Errorf("wrong validate type, got %T, expected a string", v))
	}
}

// ValidateDaysInWeekOrMonth validates the DaysInWeekOrMonth field
// The day range should be in (0,31]
func ValidateDaysInWeekOrMonth(v, _ interface{}) bool {
	switch v := v.(type) {
	case map[int]bool:
		for k := range v {
			if k <= 0 || k > 31 {
				return false
			}
		}
		return true
	default:
		panic(fmt.Errorf("wrong validate type, got %T, expected an int set (map[int]bool)", v))
	}
}

// RepeatReminder represents a reminder for habit
type RepeatReminder struct {
	Cycle             string        `bson:"cycle" json:"cycle" valid:"required,cycle"`
	DaysInWeekOrMonth map[int]bool  `bson:"days" json:"repeat_on" valid:"required,daysInWeekOrMonth"`
	TimeInDay         time.Duration `bson:"remindAt" json:"remindAt" valid:"-"`
	Duration          time.Duration `bson:"duration" json:"duration" valid:"-"`
}
