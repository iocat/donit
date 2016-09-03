package habits

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
	col = utils.MakeMGOCollectionFunc(utils.Habit)
}

// CreateDoc creates a document
func CreateDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, h *Habit) (bson.ObjectId, error) {
	l.Log("ctx", ctx, "op", "habits.CreateDoc", "Habit", h)
	// Generate a new id
	h.ObjectId = bson.NewObjectId()
	if err := utils.Validate(h); err != nil {
		l.Log("ctx", ctx, "result", fmt.Errorf("validation error: %s", err))
		return "", errors.NewValidate(err.Error())
	}
	if err := col(db).Insert(h); err != nil {
		// Do not expect id duplication, so duplicate error is not checked here
		l.Log("ctx", ctx, "result", fmt.Errorf("error: %s", err))
		return "", errors.NewDuplicated("habit", h.ObjectId.Hex())
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return h.ObjectId, nil
}

// DeleteDoc deletes a document
func DeleteDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, id bson.ObjectId) error {
	l.Log("ctx", ctx, "op", "habits.DeleteDoc", "habit", id)
	if err := col(db).RemoveId(id); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", ctx, "result", "NOT_FOUND")
			return errors.NewNotFound("habit", id.Hex())
		}
		l.Log("ctx", ctx, "result", fmt.Errorf("error: %s", err))
		return err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return nil
}

// UpdateDoc updates a document
func UpdateDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, h *Habit) error {
	l.Log("ctx", ctx, "op", "habits.UpdateDoc", "habit", h)
	if err := utils.Validate(h); err != nil {
		l.Log("ctx", ctx, "result", fmt.Errorf("validation error: %s", err))
		return errors.NewValidate(err.Error())
	}
	if err := col(db).UpdateId(h.ObjectId, h); err != nil {
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
func ReadDoc(l log.Logger, ctx utils.UUID, db *mgo.Database, id bson.ObjectId) (*Habit, error) {
	l.Log("ctx", ctx, "op", "habits.ReadDoc", "habit", id)
	var h Habit
	if err := col(db).FindId(id).One(&h); err != nil {
		if err == mgo.ErrNotFound {
			l.Log("ctx", ctx, "result", "NOT_FOUND")
			return nil, errors.NewNotFound("habit", id.Hex())
		}
		l.Log("ctx", ctx, "result", err)
		return nil, err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return &h, nil
}
