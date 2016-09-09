package goal

import (
	"fmt"

	"github.com/iocat/donit/internal/achieving/errors"
	"github.com/iocat/donit/internal/achieving/internal/achievable"
	"github.com/iocat/donit/internal/achieving/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// AddTask adds a habit
func (g *Goal) AddTask(ac *mgo.Collection, a *achievable.Habit) (utils.HexID, error) {
	id := utils.HexID{ObjectId: bson.NewObjectId()}
	a.Goal, a.HexID = g.HexID, id
	err := ac.Insert(a)
	if err != nil {
		// Does not catch duplication error
		return id, err
	}
	return id, nil
}

// RemoveTask removes a habit
func (g *Goal) RemoveTask(ac *mgo.Collection, id utils.HexID) error {
	err := ac.Remove(bson.M{
		"goal": g.HexID,
		"_id":  id,
	})
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.NewNotFound("task", fmt.Sprintf("%s,%s", g.HexID, id))
		}
	}
	return nil
}

// UpdateTask updates a task
func (g *Goal) UpdateTask(ac *mgo.Collection, a *achievable.Task, id utils.HexID) error {
	err := ac.Update(bson.M{
		"goal": g.HexID,
		"_id":  id,
	}, a)
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.NewNotFound("task", fmt.Sprintf("%s,%s", g.HexID, id))
		}
	}
	return nil
}

// RetrieveTask gets the task list
func (g *Goal) RetrieveTask(a *mgo.Collection, limit, offset int) ([]achievable.Task, error) {
	var h []achievable.Task
	err := g.retrieve(&h, a, limit, offset)
	if err != nil {
		return nil, err
	}
	return h, nil
}
