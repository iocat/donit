package data

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Comment represents a comment
type Comment struct {
	User string        `bson:"user" json:"-"`
	Goal bson.ObjectId `bson:"goal" json:"-"`
	ID   bson.ObjectId `bson:"_id,omitempty" json:"id"`

	FromUser string    `bson:"from" json:"from"`
	Comment  string    `bson:"comment" json:"comment"`
	At       time.Time `bson:"lastUpdated" json:"lastUpdated"`
	Edited   *bool     `bson:"edited" json:"edited"`
}

func (c *Comment) cname() string {
	return CollectionComment
}

func (c *Comment) keys() bson.M {
	return bson.M{
		"user": c.User,
		"goal": c.Goal,
		"_id":  c.ID,
	}
}

func (c *Comment) SetKeys(k []string) error {
	if err := writeKey(k, &c.User, &c.Goal, &c.ID); err != nil {
		return err
	}
	return nil
}

func (c *Comment) Validate() error {
	switch {
	case len(c.FromUser) == 0:
		return newBadData("comment from unknown user, field from not provided")
	case len(c.Comment) == 0:
		return newBadData("comment length is 0")
	case c.Edited == nil:
		return newBadData("the edited field is not provided")
	}
	c.At = time.Now()
	return nil
}
