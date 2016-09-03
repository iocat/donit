package tasks

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
	col = utils.MakeMGOCollectionFunc(utils.Task)
}

// CreateDoc creates a document
func CreateDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, t *Task) (bson.ObjectId, error) {
	l.Log("ctx", ctx, "op", "tasks.CreateDoc", "task", t)
	// Generate a new id
	t.ObjectId = bson.NewObjectId()
	if err := utils.Validate(t); err != nil {
		l.Log("ctx", ctx, "result", fmt.Errorf("validation error: %s", err))
		return "", errors.NewValidate(err.Error())
	}
	if err := col(db).Insert(t); err != nil {
		// Do not expect id duplication, so duplicate error is not checked here
		l.Log("ctx", ctx, "result", fmt.Errorf("error: %s", err))
		return "", errors.NewDuplicated("task", t.ObjectId.Hex())
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return t.ObjectId, nil
}

// DeleteDoc deletes a document
func DeleteDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, id bson.ObjectId) error {
	l.Log("ctx", ctx, "op", "tasks.DeleteDoc", "task", id)
	if err := col(db).RemoveId(id); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", ctx, "result", "NOT_FOUND")
			return errors.NewNotFound("task", id.Hex())
		}
		l.Log("ctx", ctx, "result", fmt.Errorf("error: %s", err))
		return err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return nil
}

// UpdateDoc updates a document
func UpdateDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, t *Task) error {
	l.Log("ctx", ctx, "op", "tasks.UpdateDoc", "task", t)
	if err := utils.Validate(t); err != nil {
		l.Log("ctx", ctx, "result", fmt.Errorf("validation error: %s", err))
		return errors.NewValidate(err.Error())
	}
	if err := col(db).UpdateId(t.ObjectId, t); err != nil {
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

// ReadDoc reads a document
func ReadDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, id bson.ObjectId) (*Task, error) {
	l.Log("ctx", ctx, "op", "tasks.ReadDoc", "task", id)
	var t Task
	if err := col(db).FindId(id).One(&t); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", ctx, "result", "NOT_FOUND")
			return nil, errors.NewNotFound("task", id.Hex())
		}
		l.Log("ctx", ctx, "result", err)
		return nil, err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return &t, nil
}
