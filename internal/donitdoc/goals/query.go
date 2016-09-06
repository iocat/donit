package goals

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/errors"

	"github.com/iocat/donit/internal/donitdoc/internal/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var col func(*mgo.Database) *mgo.Collection

func init() {
	col = utils.Goal.Collection
}

// CreateDoc creates a new goal and returns a new id generated regardless of
// the id inside the provided Goal object
func CreateDoc(l log.Logger, ctx string, db *mgo.Database, g *Goal) (bson.ObjectId, error) {
	l.Log("ctx", ctx, "op", "goals.CreateDoc", "goal", g)
	// Generate a new id
	g.ObjectId = bson.NewObjectId()
	if err := utils.Validate(g); err != nil {
		l.Log("ctx", ctx, "result", fmt.Errorf("validation error: %s", err))
		return "", errors.NewValidate(err.Error())
	}
	if err := col(db).Insert(g); err != nil {
		// Do not expect id duplication, so duplicate error is not checked here
		l.Log("ctx", ctx, "result", fmt.Errorf("error: %s", err))
		return "", errors.NewDuplicated("goal", g.ObjectId.Hex())
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return g.ObjectId, nil
}

// DeleteDoc deletes a goal
func DeleteDoc(l log.Logger, ctx string, db *mgo.Database, id bson.ObjectId) error {
	l.Log("ctx", ctx, "op", "goals.DeleteDoc", "goal", id)
	if err := col(db).RemoveId(id); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", ctx, "result", "NOT_FOUND")
			return errors.NewNotFound("goal", id.Hex())
		}
		l.Log("ctx", ctx, "result", fmt.Errorf("error: %s", err))
		return err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return nil
}

// UpdateDoc updates a goal
func UpdateDoc(l log.Logger, ctx string, db *mgo.Database, g *Goal) error {
	l.Log("ctx", ctx, "op", "goals.UpdateDoc", "goal", g)
	if err := utils.Validate(g); err != nil {
		l.Log("ctx", ctx, "result", fmt.Errorf("validation error: %s", err))
		return errors.NewValidate(err.Error())
	}
	if err := col(db).UpdateId(g.ObjectId, g); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", ctx, "result", "NOT_FOUND")
			return err
		}
		l.Log("ctx", ctx, "result", err)
		return err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return nil
}

// ReadDoc reads a goal from the database
func ReadDoc(l log.Logger, ctx string, db *mgo.Database, id bson.ObjectId) (*Goal, error) {
	l.Log("ctx", ctx, "op", "goals.ReadDoc", "goal", id)
	var g Goal
	if err := col(db).FindId(id).One(&g); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", ctx, "result", "NOT_FOUND")
			return nil, errors.NewNotFound("goal", id.Hex())
		}
		l.Log("ctx", ctx, "result", err)
		return nil, err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return &g, nil
}

// AllDocsOfUser gets all the goals associated with this user
func AllDocsOfUser(l log.Logger, ctx string, db *mgo.Database, username string, limit, offset int) ([]Goal, error) {
	l.Log("ctx", ctx, "op", "tasks.AllDocsOfUser", "username", username, "limit", limit, "offset", offset)
	var h []Goal
	q := utils.Query{
		Query: col(db).Find(bson.M{
			"username": username,
		}),
	}
	if err := q.Limit(limit).Skip(offset).All(&h); err != nil {
		l.Log("ctx", ctx, "result", err)
		return nil, nil
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return h, nil
}

// DeleteAllDocsOfUser deletes all the goals of this user (without the linked goals)
func DeleteAllDocsOfUser(l log.Logger, ctx string, db *mgo.Database, username string) error {
	l.Log("ctx", ctx, "op", "goals.DeleteAllDocsOfUser", "username", username)
	// Remove all the goals
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
