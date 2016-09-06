package habits

import (
	"fmt"
	"time"

	"github.com/iocat/donit/internal/donitdoc/achievable"
	valid "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

const (
	// EveryDay is the Cycle for Every day habit
	EveryDay = "EVERYDAY"
	// EveryWeekAndCustom is the cycle for every week habit
	EveryWeekAndCustom = "EVERY_WEEK"
	// EveryMonthAndCustom is the cycle for every month habit
	EveryMonthAndCustom = "EVERY_MONTH"
)
const (
	// Sunday represents Sunday
	Sunday = iota + 1
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

func init() {
	valid.SetFieldsRequiredByDefault(true)
	valid.CustomTypeTagMap.Set("cycle", valid.CustomTypeValidator(validateCycle))
	valid.CustomTypeTagMap.Set("daysInWeekOrMonth", valid.CustomTypeValidator(validateDaysInWeekOrMonth))
	valid.CustomTypeTagMap.Set("validateStatus", valid.CustomTypeValidator(achievable.ValidateStatus))
}

// validateCycle validates the Cycle field
func validateCycle(v, _ interface{}) bool {
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

// validateDaysInWeekOrMonth validates the DaysInWeekOrMonth field
// The day range should be in (0,31]
func validateDaysInWeekOrMonth(v, _ interface{}) bool {
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

// Habit represents a goal's habit
type Habit struct {
	bson.ObjectId         `bson:"_id,omitempty" json:"id" valid:"required,hexadecimal"`
	Goal                  bson.ObjectId `bson:"goal" json:"-" valid:"required,hexadecimal"`
	achievable.Achievable `bson:"subGoal,inline" valid:"required"`
	*RepeatReminder       `bson:"reminder,omitempty" json:"reminder,omitempty" valid:"optional"`
}
