// Copyright 2016 Thanh Ngo <felix.infinite@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package goal contains goal data
package goal

import (
	"fmt"
	"time"

	"github.com/iocat/donit/internal/achieving/errors"
	"github.com/iocat/donit/internal/achieving/internal/achievable"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// VisiblePrivate only allows the user to see
	VisiblePrivate = "PRIVATE"
	// VisibleForFollowers allows followers to see only
	VisibleForFollowers = "FOR_FOLLOWERS"
	// VisiblePublic allows everyone to see it
	VisiblePublic = "PUBLIC"
)

// Goal represents an achievable Goal
type Goal struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty" valid:"optional,hexadecimal"`
	Username string        `bson:"username" json:"-" valid:"optional,alphanum,length(1|30)"`

	Name        string                  `bson:"name" json:"name" valid:"name" valid:"required,utfletternum,stringlength(1|100)"`
	Description string                  `bson:"description,omitempty" json:"description,omitempty" valid:"optional,stringlength(1|400)"`
	LastUpdated time.Time               `bson:"lastUpdated" json:"lastUpdated" valid:"-"`
	Status      string                  `bson:"status" json:"status" valid:"validateStatus"`
	PictureURL  string                  `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty" valid:"optional,url"`
	Visibility  string                  `bson:"visibility" json:"visibility,omitempty" valid:"required,goalVisibilityValidator"`
	ToDo        []achievable.Achievable `bson:"-" json:"achievables" valid:"-"`
}

// VisibilityValidatorFunc validates the visibility field of goal
func VisibilityValidatorFunc(value, _ interface{}) bool {
	switch value := value.(type) {
	case string:
		switch value {
		case VisiblePublic, VisiblePrivate, VisibleForFollowers:
			return true
		default:
			return false
		}
	default:
		panic("the visibility field must be a string")
	}
}

// RemoveAchievable removes a habit
func (g *Goal) RemoveAchievable(ac *mgo.Collection, id bson.ObjectId) error {
	err := ac.Remove(bson.M{
		"_goal": g.ID,
		"_id":   id,
	})
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.NewNotFound("achievable", fmt.Sprintf("%s,%s", g.ID.Hex(), id))
		}
	}
	return nil
}

// UpdateAchievable updates an achievable task
func (g *Goal) UpdateAchievable(ac *mgo.Collection, a *achievable.Achievable, id bson.ObjectId) error {
	a.Goal, a.ID = g.ID, id
	err := ac.Update(bson.M{
		"_goal": g.ID,
		"_id":   id,
	}, a)
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.NewNotFound("habit", fmt.Sprintf("%s,%s", g.ID, id))
		}
	}
	return nil
}

func (g *Goal) retrieveAchievables(list *[]achievable.Achievable, a *mgo.Collection, limit, offset int) error {
	q := a.Find(bson.M{
		"_goal": g.ID,
	})
	if limit > 0 {
		q.Limit(limit)
	}
	if offset > 0 {
		q.Skip(offset)
	}
	err := q.All(list)
	if err != nil {
		return err
	}
	return nil
}

// RetrieveAchievables gets the habit list
func (g *Goal) RetrieveAchievables(ac *mgo.Collection, limit, offset int) ([]achievable.Achievable, error) {
	var h []achievable.Achievable
	err := g.retrieveAchievables(&h, ac, limit, offset)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// AddAchievable adds a habit
func (g *Goal) AddAchievable(ac *mgo.Collection, a *achievable.Achievable) (bson.ObjectId, error) {
	id := bson.NewObjectId()
	a.Goal, a.ID = g.ID, id
	err := ac.Insert(a)
	if err != nil {
		// Does not catch duplication error
		return id, err
	}
	return id, nil
}
