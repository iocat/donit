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
}
