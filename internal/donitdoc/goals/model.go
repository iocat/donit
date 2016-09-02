package goals

import (
	"github.com/iocat/donit/internal/donitdoc/achievable"
	valid "gopkg.in/asaskevich/govalidator.v4"
)

const (
	// AccessPrivate is private accessibility
	AccessPrivate = "PRIVATE"
	// AccessForFollowers is the accessibility for followers
	AccessForFollowers = "FOR_FOLLOWERS"
	// AccessPublic is the accessibility for public user
	AccessPublic = "PUBLIC"
)

func init() {
	valid.SetFieldsRequiredByDefault(true)
	valid.CustomTypeTagMap.Set("goalAccessValidator", valid.CustomTypeValidator(GoalAccessValidatorFunc))
}

// Goal represents an achievable Goal
type Goal struct {
	achievable.Achievable `bson:"subGoal,inline" valid:"required"`
	PictureURL            *string `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty" valid:"optional,url"`
	Accessibility         string  `bson:"accessibility" json:"accessibility,omitempty" valid:"required,goalAccessValidator"`
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
