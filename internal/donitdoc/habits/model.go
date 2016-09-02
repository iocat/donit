package habits

import (
	"time"

	"github.com/iocat/donit/internal/donitdoc/achievable"
	"gopkg.in/mgo.v2/bson"
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
