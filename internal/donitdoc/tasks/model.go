package tasks

import (
	"time"

	"github.com/iocat/donit/internal/donitdoc/achievable"

	"gopkg.in/mgo.v2/bson"
)

// Reminder represents a reminder for tasks
type Reminder struct {
	At       time.Time     `bson:"remindAt" json:"remindAt"`
	Duration time.Duration `bson:"duration" json:"duration"`
}

// Task represents a task
type Task struct {
	bson.ObjectId         `bson:"_id,omitempty" json:"id"`
	achievable.Achievable `bson:"subGoal,inline"`
	*Reminder             `bson:"reminder,omitempty" json:"reminder,omitempty"`
}
