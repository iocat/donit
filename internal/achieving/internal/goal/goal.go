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

	"github.com/iocat/donit/internal/achieving/errors"
	"github.com/iocat/donit/internal/achieving/internal/achievable"
	"github.com/iocat/donit/internal/achieving/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// AccessPrivate is private accessibility
	AccessPrivate = "PRIVATE"
	// AccessForFollowers is the accessibility for followers
	AccessForFollowers = "FOR_FOLLOWERS"
	// AccessPublic is the accessibility for public user
	AccessPublic = "PUBLIC"
)

// Goal represents an achievable Goal
type Goal struct {
	utils.HexID           `bson:"id,inline" valid:"required"`
	Username              string `bson:"username" json:"username" valid:"required,alphanum,length(1|30)"`
	achievable.Achievable `bson:"subGoal,inline" valid:"required"`
	PictureURL            string                  `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty" valid:"optional,url"`
	Accessibility         string                  `bson:"accessibility" json:"accessibility,omitempty" valid:"required,goalAccessValidator"`
	ToDo                  []achievable.Achievable `bson:"-" json:"todo" valid:"-"`
}

// AccessValidatorFunc validates the accessibility field of the Goal model
func AccessValidatorFunc(value, _ interface{}) bool {
	switch value := value.(type) {
	case string:
		switch value {
		case AccessPrivate, AccessPublic, AccessForFollowers:
			return true
		default:
			return false
		}
	default:
		panic("the accessibility field must be a string")
	}
}

// RetrieveHabit gets the habit list
func (g *Goal) retrieve(list *[]achievable.Achievable, a *mgo.Collection, limit, offset int) error {
	q := a.Find(bson.M{
		"goal": g.HexID,
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

// RemoveAchievable removes a habit
func (g *Goal) RemoveAchievable(ac *mgo.Collection, id utils.HexID) error {
	err := ac.Remove(bson.M{
		"goal": g.HexID,
		"_id":  id,
	})
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.NewNotFound("achievable", fmt.Sprintf("%s,%s", g.HexID, id))
		}
	}
	return nil
}

// UpdateAchievable updates an achievable task
func (g *Goal) UpdateAchievable(ac *mgo.Collection, a *achievable.Achievable, id utils.HexID) error {
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

// RetrieveAchievable gets the habit list
func (g *Goal) RetrieveAchievable(ac *mgo.Collection, limit, offset int) ([]achievable.Achievable, error) {
	var h []achievable.Achievable
	err := g.retrieve(&h, ac, limit, offset)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// AddAchievable adds a habit
func (g *Goal) AddAchievable(ac *mgo.Collection, a *achievable.Achievable) (utils.HexID, error) {
	id := utils.HexID{ObjectId: bson.NewObjectId()}
	a.Goal, a.HexID = g.HexID, id
	err := ac.Insert(a)
	if err != nil {
		// Does not catch duplication error
		return id, err
	}
	return id, nil
}
