package comments

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/errors"
	"github.com/iocat/donit/internal/donitdoc/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var col func(*mgo.Database) *mgo.Collection

func init() {
	col = utils.MakeMGOCollectionFunc(utils.Comment)
}

// CreateDoc creates a new goal and returns a new id generated regardless of
// the id inside the provided Goal object
func CreateDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, c *Comment) (bson.ObjectId, error) {
	l.Log("ctx", ctx, "op", "comments.CreateDoc", "comment", c)
	// Generate a new id
	c.ObjectId = bson.NewObjectId()
	if err := utils.Validate(c); err != nil {
		l.Log("ctx", ctx, "result", fmt.Errorf("validation error: %s", err))
		return "", errors.NewValidate(err.Error())
	}
	if err := col(db).Insert(c); err != nil {
		// Do not expect id duplication, so duplicate error is not checked here
		l.Log("ctx", ctx, "result", fmt.Errorf("error: %s", err))
		return "", errors.NewDuplicated("comment", c.ObjectId.Hex())
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return c.ObjectId, nil
}

// DeleteDoc deletes a goal
func DeleteDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, id bson.ObjectId) error {
	l.Log("ctx", ctx, "op", "comments.DeleteDoc", "goal", id)
	if err := col(db).RemoveId(id); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", ctx, "result", "NOT_FOUND")
			return errors.NewNotFound("comment", id.Hex())
		}
		l.Log("ctx", ctx, "result", fmt.Errorf("error: %s", err))
		return err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return nil
}

// UpdateDoc updates a goal
func UpdateDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, c *Comment) error {
	l.Log("ctx", ctx, "op", "comments.UpdateDoc", "comment", c)
	if err := utils.Validate(c); err != nil {
		l.Log("ctx", ctx, "result", fmt.Errorf("validation error: %s", err))
		return errors.NewValidate(err.Error())
	}
	if err := col(db).UpdateId(c.ObjectId, c); err != nil {
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
func ReadDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, id bson.ObjectId) (*Comment, error) {
	l.Log("ctx", ctx, "op", "comments.ReadDoc", "comment", id)
	var c Comment
	if err := col(db).FindId(id).One(&c); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", ctx, "result", "NOT_FOUND")
			return nil, errors.NewNotFound("comment", id.Hex())
		}
		l.Log("ctx", ctx, "result", err)
		return nil, err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return &c, nil
}
