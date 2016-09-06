package achievable

import (
	"time"
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

const (
	// Done represents achievable done status
	Done = "DONE"
	// NotDone represents achievable not done status
	NotDone = "NOT_DONE"
	// InProgress represents achievable InProgress status
	InProgress = "IN_PROGRESS"
)

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
