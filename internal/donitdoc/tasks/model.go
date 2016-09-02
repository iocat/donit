package tasks

import (
	"time"

	"github.com/iocat/donit/internal/donitdoc/achievable"
	valid "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	valid.SetFieldsRequiredByDefault(true)
}

// Reminder represents a reminder for tasks
type Reminder struct {
	At       time.Time     `bson:"remindAt" json:"remindAt" valid:"-"`
	Duration time.Duration `bson:"duration" json:"duration" valid:"-"`
}

// Task represents a task
type Task struct {
	bson.ObjectId         `bson:"_id,omitempty" json:"id" valid:"required,hexadecimal"`
	achievable.Achievable `bson:"subGoal,inline" valid:"required"`
	*Reminder             `bson:"reminder,omitempty" json:"reminder,omitempty" valid:"optional"`
}
