package habits

import (
	"time"

	"github.com/iocat/donit/internal/donitdoc/achievable"
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

// RepeatReminder represents a reminder for habit
type RepeatReminder struct {
	Cycle             string        `bson:"cycle" json:"cycle"`
	DaysInWeekOrMonth map[int]bool  `bson:"days" json:"repeat_on"`
	TimeInDay         time.Duration `bson:"remindAt" json:"remindAt"`
	Duration          time.Duration `bson:"duration" json:"duration"`
}

// Habit represents a goal's habit
type Habit struct {
	bson.ObjectId         `bson:"_id,omitempty" json:"id"`
	achievable.Achievable `bson:"subGoal,inline"`
	*RepeatReminder       `bson:"reminder,omitempty" json:"reminder,omitempty"`
}
