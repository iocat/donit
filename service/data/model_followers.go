package data

import "gopkg.in/mgo.v2/bson"

// Follower represents the user follower
type Follower struct {
	User     string `bson:"username" json:"-"`
	Follower string `bson:"follower" json:"follower"`
}

// cname implements the Item's cname
func (f *Follower) cname() string {
	return CollectionFollower
}

// keys implements the Item's keys
func (f *Follower) keys() bson.M {
	return bson.M{
		"user":     f.User,
		"follower": f.Follower,
	}
}

// SetKeys implements Item's SetKeys
func (f *Follower) SetKeys(k []string) error {
	if err := writeKey(k, &f.User, &f.Follower); err != nil {
		return err
	}
	return nil
}
