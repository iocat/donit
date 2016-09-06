package users

import (
	"time"

	"github.com/iocat/donit/internal/donitdoc/goals"
	valid "gopkg.in/asaskevich/govalidator.v4"
)

const (
	// Offline represents user offline status
	Offline = "OFFLINE"
	// OnlineAvailable represents user online status
	OnlineAvailable = "ONLINE"
	// Busy represents user busy status
	Busy = "BUSY"
)

func validateUserStatus(value, _ interface{}) bool {
	switch value := value.(type) {
	case string:
		switch value {
		case Offline, OnlineAvailable, Busy:
			return true
		default:
			return false
		}
	default:
		panic("the accessibility field must be a string")
	}
}

func init() {
	valid.SetFieldsRequiredByDefault(true)
	valid.CustomTypeTagMap.Set("goalAccessField", valid.CustomTypeValidator(goals.GoalAccessValidatorFunc))
	valid.CustomTypeTagMap.Set("validateUserStatus", valid.CustomTypeValidator(validateUserStatus))
}

// User represents a user
type User struct {
	Username             string    `bson:"username" json:"username" valid:"required,alphanum,length(1|30)"`
	Status               string    `bson:"status" json:"status" valid:"required,validateUserStatus"`
	Email                string    `bson:"email" json:"email" valid:"required,email"`
	Firstname            string    `bson:"firstName" json:"firstName" valid:"required,alpha,length(0|50)"`
	Lastname             string    `bson:"lastName" json:"lastName" valid:"required,alpha,length(0|50)"`
	DefaultAccessibility string    `bson:"defaultAccess" json:"defaultAccess" valid:"required,goalAccessField"`
	PictureURL           *string   `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty" valid:"optional,url"`
	LastUpdated          time.Time `bson:"lastUpdated" json:"lastUpdated" valid:"-"`
	HasUpdate            bool      `bson:"hasUpdated" json:"hasUpdated" valid:"-"`
}

// StoredUser encapsulates user's password
type StoredUser struct {
	User     `valid:"required"`
	Password string `bson:"password" json:"password,omitempty" valid:"hexadecimal,optional"`
	Salt     string `bson:"salt" json:"-" valid:"alphanum,optional"`
}
