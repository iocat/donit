package followers

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/errors"
	"github.com/iocat/donit/internal/donitdoc/internal/utils"
	valid "gopkg.in/asaskevich/govalidator.v4"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// CollectionName for secondary indexing purpose
var CollectionName = utils.Follower.MGOName()

var col func(*mgo.Database) *mgo.Collection

func init() {
	valid.SetFieldsRequiredByDefault(true)
	col = utils.Follower.Collection
}

// CreateDoc creates a new follower
func CreateDoc(l log.Logger, ctx string, db *mgo.Database, g *Follower) error {
	l.Log("ctx", ctx, "op", "followers.CreateDoc", "goal", g)
	if err := utils.Validate(g); err != nil {
		l.Log("ctx", ctx, "result", err)
		return errors.NewValidate(err.Error())
	}
	if err := col(db).Insert(g); err != nil {
		// Do not expect id duplication, so duplicate error is not checked here
		l.Log("ctx", ctx, "result", fmt.Errorf("error: %s", err))
		return errors.NewDuplicated("follower", fmt.Sprintf("(%s,%s)", g.Username, g.Follower))
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return nil
}

// DeleteDoc deletes a follower
func DeleteDoc(l log.Logger, ctx string, db *mgo.Database, username string, follower string) error {
	l.Log("ctx", ctx, "op", "followers.DeleteDoc", "username", username, "follower", follower)
	if err := col(db).Remove(bson.M{
		"username": username,
		"follower": follower,
	}); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", ctx, "result", "NOT_FOUND")
			return errors.NewNotFound("follower", fmt.Sprintf("(%s,%s)", username, follower))
		}
		l.Log("ctx", ctx, "result", fmt.Errorf("error: %s", err))
		return err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return nil
}

// AllDocsOfUser gets all the followers associated with this user
func AllDocsOfUser(l log.Logger, ctx string, db *mgo.Database, username string, limit, offset int) ([]Follower, error) {
	l.Log("ctx", ctx, "op", "followers.AllDocsOfUser", "username", username, "limit", limit, "offset", offset)
	var h []Follower
	q := &utils.Query{Query: col(db).Find(bson.M{"username": username})}
	if err := q.Limit(limit).Skip(offset).All(&h); err != nil {
		l.Log("ctx", ctx, "result", err)
		return nil, err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return h, nil
}

// DeleteAllDocsOfUser deletes all the followers of this user
func DeleteAllDocsOfUser(l log.Logger, ctx string, db *mgo.Database, username string) error {
	l.Log("ctx", ctx, "op", "followers.DeleteAllDocsOfUser", "username", username)
	if _, err := col(db).RemoveAll(
		bson.M{
			"username": username,
		},
	); err != nil {
		l.Log("ctx", ctx, "result", err)
		return err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return nil
}
