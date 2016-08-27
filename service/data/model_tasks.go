package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Reminder represents a reminder for tasks
type Reminder struct {
	At       time.Time     `bson:"remindAt" json:"remindAt"`
	Duration time.Duration `bson:"duration" json:"duration"`
}

// Task represents a task
type Task struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id"`
	SubGoal   `bson:"subGoal,inline"`
	User      string        `bson:"user" json:"-"`
	Goal      bson.ObjectId `bson:"goal" json:"-"`
	*Reminder `bson:"reminder,omitempty" json:"reminder,omitempty"`
}

// Cname implements Item's cname
func (t *Task) cname() string {
	return CollectionTask
}

// KeySet implements Item's keys
func (t *Task) keys() bson.M {
	return bson.M{
		"_id":  t.ID,
		"user": t.User,
		"goal": t.Goal,
	}
}

// SetKeys implements Item's SetKeys
func (t *Task) SetKeys(k []string) error {
	if err := writeKey(k, &t.User, &t.Goal, &t.ID); err != nil {
		return err
	}
	return nil
}
