package achievable

import "time"

// Achievable represents an achievable action
// Achievable IS NOT a document in database, it is meant to be embedded
// inside other documents
type Achievable struct {
	Name        string    `bson:"name" json:"name"`
	Description *string   `bson:"description,omitempty" json:"description,omitempty"`
	LastUpdated time.Time `bson:"createdAt" json:"createdAt"`
}
