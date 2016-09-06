package comments

import (
	"time"

	valid "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2/bson"
)

func init() {
	valid.SetFieldsRequiredByDefault(true)
}

// Comment represents a comment
type Comment struct {
	bson.ObjectId `bson:"_id,omitempty" json:"id" valid:"required,hexadecimal"`
	Who           string    `bson:"username" json:"username" valid:"required,alphanum,length(1|30)"`
	Content       string    `bson:"content" json:"content" valid:"required,utfletternum,stringlength(0|1000)"`
	At            time.Time `bson:"lastUpdated" json:"lastUpdated" valid:"-"`
	Edited        bool      `bson:"edited" json:"edited" valid:"-"`
}
