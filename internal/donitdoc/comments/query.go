package comments

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
	col = utils.Comment.Collection
}

// CreateDoc creates a new goal and returns a new id generated regardless of
// the id inside the provided Goal object
func CreateDoc(l log.Logger, ctx string, db *mgo.Database, c *Comment) (bson.ObjectId, error) {
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
func DeleteDoc(l log.Logger, ctx string, db *mgo.Database, id bson.ObjectId) error {
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
func UpdateDoc(l log.Logger, ctx string, db *mgo.Database, c *Comment) error {
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
func ReadDoc(l log.Logger, ctx string, db *mgo.Database, id bson.ObjectId) (*Comment, error) {
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

// ALlDocsOfGoal gets all the goals associated with this goal
func ALlDocsOfGoal(l log.Logger, ctx string, db *mgo.Database, goal bson.ObjectId, limit, offset int) ([]Comment, error) {
	l.Log("ctx", ctx, "op", "tasks.AllDocsOfUser", "goal", goal, "limit", limit, "offset", offset)
	var h []Comment
	q := utils.Query{
		Query: col(db).Find(bson.M{
			"goal": goal,
		}),
	}
	if err := q.Limit(limit).Skip(offset).All(&h); err != nil {
		l.Log("ctx", ctx, "result", err)
		return nil, nil
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return h, nil
}

// DeleteAllDocsOfGoal deletes all comments of the goal
func DeleteAllDocsOfGoal(l log.Logger, ctx string, db *mgo.Database, goal bson.ObjectId) error {
	l.Log("ctx", ctx, "op", "comments.DeleteAllDocsOfGoal", "goal", goal)
	if _, err := col(db).RemoveAll(
		bson.M{
			"goal": goal,
		},
	); err != nil {
		l.Log("ctx", ctx, "result", err)
		return err
	}
	l.Log("ctx", ctx, "result", "SUCCESS")
	return nil
}
