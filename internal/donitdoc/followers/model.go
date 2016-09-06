package followers

import "time"

// Follower represents a follower
type Follower struct {
	Username string    `bson:"username" json:"username" valid:"required,alphanum,length(1|30)"`
	Follower string    `bson:"follower" json:"username" valid:"required,alphanum,length(1|30)"`
	FollowAt time.Time `bson:"followAt" json:"followAt" valid:"-"`
}
