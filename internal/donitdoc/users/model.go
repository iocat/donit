package users

import (
	"fmt"
	"time"

	"github.com/iocat/donit/internal/donitdoc/goals"
	valid "gopkg.in/asaskevich/govalidator.v4"
)

func init() {
	valid.SetFieldsRequiredByDefault(true)
	valid.CustomTypeTagMap.Set("userStatusField", valid.CustomTypeValidator(validateUserStatusField))
	valid.CustomTypeTagMap.Set("goalAccessField", valid.CustomTypeValidator(goals.GoalAccessValidatorFunc))
}

const (
	statusOffline         = "OFFLINE"
	statusOnlineAvailable = "ONLINE"
	statusBusy            = "BUSY"
)

func validateUserStatusField(value, _ interface{}) bool {
	switch value := value.(type) {
	case string:
		switch value {
		case statusOffline, statusOnlineAvailable, statusBusy:
			return true
		default:
			return false
		}
	default:
		panic(fmt.Errorf("user status field must be a string, got %T", value))
	}
}

// User represents a user
type User struct {
	Username             string    `bson:"username" json:"username" valid:"required,alphanum,length(1|30)"`
	Email                string    `bson:"email" json:"email" valid:"required,email"`
	Firstname            string    `bson:"firstName" json:"firstName" valid:"required,alpha,length(0|50)"`
	Lastname             string    `bson:"lastName" json:"lastName" valid:"required,alpha,length(0|50)"`
	Status               string    `bson:"status" json:"status" valid:"userStatusField" `
	DefaultAccessibility string    `bson:"defaultAccess" json:"defaultAccess" valid:"goalAccessField"`
	LastUpdated          time.Time `bson:"lastUpdated" json:"lastUpdated" valid:"-"`
	HasUpdate            bool      `bson:"hasUpdated" json:"hasUpdated" valid:"-"`
	PictureURL           *string   `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty" valid:"URL,optional"`
}

// StoredUser encapsulates user's password
type StoredUser struct {
	User     `valid:"required"`
	Password *string `bson:"password" json:"password,omitempty" valid:"hexadecimal,optional"`
	Salt     *string `bson:"salt" json:"-" valid:"alphanum,optional"`
}
