package goals

import (
	"github.com/iocat/donit/internal/donitdoc/achievable"
)

const (
	// AccessPrivate is private accessibility
	AccessPrivate = "PRIVATE"
	// AccessForFollowers is the accessibility for followers
	AccessForFollowers = "FOR_FOLLOWERS"
	// AccessPublic is the accessibility for public user
	AccessPublic = "PUBLIC"
)

// Goal represents an achievable Goal
type Goal struct {
	achievable.Achievable `bson:"subGoal,inline"`
	PictureURL            *string `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty"`
	Accessibility         string  `bson:"accessibility" json:"accessibility,omitempty"`
}

// GoalAccessValidatorFunc validates the accessibility field of the Goal model
func GoalAccessValidatorFunc(value, _ interface{}) bool {
	switch value := value.(type) {
	case string:
		switch value {
		case AccessPrivate, AccessPublic, AccessForFollowers:
			return true
		default:
			return false
		}
	default:
		panic("the accessibility field must be a string")
	}
}
