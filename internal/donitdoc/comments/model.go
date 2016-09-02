package comments

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Comment represents a comment
type Comment struct {
	bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Content       string    `bson:"content" json:"content"`
	At            time.Time `bson:"lastUpdated" json:"lastUpdated"`
	Edited        *bool     `bson:"edited" json:"edited"`
}
