package goal

import (
	"fmt"

	"github.com/iocat/donit/internal/achieving/errors"
	"github.com/iocat/donit/internal/achieving/internal/achievable"
	"github.com/iocat/donit/internal/achieving/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RetrieveHabit gets the habit list
func (g *Goal) RetrieveHabit(a *mgo.Collection, limit, offset int) ([]achievable.Habit, error) {
	var h []achievable.Habit
	err := g.retrieve(&h, a, limit, offset)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// AddHabit adds a habit
func (g *Goal) AddHabit(ac *mgo.Collection, a *achievable.Habit) (utils.HexID, error) {
	id := utils.HexID{ObjectId: bson.NewObjectId()}
	a.Goal, a.HexID = g.HexID, id
	err := ac.Insert(a)
	if err != nil {
		// Does not catch duplication error
		return id, err
	}
	return id, nil
}

// RemoveHabit removes a habit
func (g *Goal) RemoveHabit(ac *mgo.Collection, id utils.HexID) error {
	err := ac.Remove(bson.M{
		"goal": g.HexID,
		"_id":  id,
	})
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.NewNotFound("habit", fmt.Sprintf("%s,%s", g.HexID, id))
		}
	}
	return nil
}

// UpdateHabit updates a habit
func (g *Goal) UpdateHabit(ac *mgo.Collection, a *achievable.Habit, id utils.HexID) error {
	err := ac.Update(bson.M{
		"goal": g.HexID,
		"_id":  id,
	}, a)
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.NewNotFound("habit", fmt.Sprintf("%s,%s", g.HexID, id))
		}
	}
	return nil
}
