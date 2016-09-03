package donitdoc

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type service struct {
	l  log.Logger
	db *mgo.Database
}

func getID(id string) (bson.ObjectId, error) {
	if !bson.IsObjectIdHex(id) {
		return "", errors.NewValidate(fmt.Sprintf("object id %s is invalid ", id))
	}
	return bson.ObjectIdHex(id), nil
}
