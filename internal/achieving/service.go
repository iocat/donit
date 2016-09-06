package achieving

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/errors"
	"github.com/iocat/donit/internal/utils"
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

func init() {
	logID = utils.MustGetLogID
}

var logID func(context.Context) string
