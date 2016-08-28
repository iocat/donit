package data

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2/bson"
)

const (
	repeatEveryDay            = "EVERYDAY"
	repeatEveryWeekAndCustom  = "EVERY_WEEK"
	repeatEveryMonthAndCustom = "EVERY_MONTH"
)
const (
	sunday = iota + 1
	monday
	tuesday
	wednesday
	thursday
	friday
	saturday
)

// RepeatReminder represents a reminder for habit
type RepeatReminder struct {
	Cycle             string        `bson:"cycle" json:"cycle"`
	DaysInWeekOrMonth map[int]bool  `bson:"days" json:"repeat_on"`
	TimeInDay         time.Duration `bson:"remindAt" json:"remindAt"`
	Duration          time.Duration `bson:"duration" json:"duration"`
}

func validateRepeat(repeat string) error {
	switch repeat {
	case repeatEveryDay, repeatEveryMonthAndCustom, repeatEveryWeekAndCustom:
		return nil
	default:
		return newBadData(fmt.Sprintf("repeat field \"%s\" is undefined (only %s, %s, and %s allowed)",
			repeat, repeatEveryDay, repeatEveryMonthAndCustom,
			repeatEveryWeekAndCustom))
	}
}

// Validate implements Validator's Validate
func (r *RepeatReminder) Validate() error {
	if err := validateRepeat(r.Cycle); err != nil {
		return err
	}
	return nil
}

// Habit represents a goal's habit
type Habit struct {
	ID              bson.ObjectId `bson:"_id,omitempty" json:"id"`
	SubGoal         `bson:"subGoal,inline"`
	User            string        `bson:"user" json:"-"`
	Goal            bson.ObjectId `bson:"goal" json:"-"`
	*RepeatReminder `bson:"reminder,omitempty" json:"reminder,omitempty"`
}

// Validate validates the habit
func (h *Habit) Validate() error {
	if err := h.SubGoal.Validate(); err != nil {
		return err
	}
	if err := h.RepeatReminder.Validate(); err != nil {
		return err
	}
	return nil
}

// cname implements Item's cname
func (h *Habit) cname() string {
	return CollectionHabit
}

// keys implements Item's keys
func (h *Habit) keys() bson.M {
	return bson.M{
		"_id":  h.ID,
		"user": h.User,
		"goal": h.Goal,
	}
}

// SetKeys implements Item's SetKeys
func (h *Habit) SetKeys(k []string) error {
	if err := writeKey(k, &h.User, &h.Goal, &h.ID); err != nil {
		return err
	}
	return nil
}

func (h *Habit) GenerateID() string {
	h.ID = bson.NewObjectId()
	return h.ID.Hex()
}
