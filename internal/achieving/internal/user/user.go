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

// Package user contains user data
package user

import (
	"fmt"
	"time"

	"github.com/iocat/donit/internal/achieving/errors"
	"github.com/iocat/donit/internal/achieving/internal/goal"
	"github.com/iocat/donit/internal/achieving/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	// Offline represents user offline status
	Offline = "OFFLINE"
	// OnlineAvailable represents user online status
	OnlineAvailable = "ONLINE"
	// Busy represents user busy status
	Busy = "BUSY"
)

// ValidateUserStatus validates the user's status field
func ValidateUserStatus(value, _ interface{}) bool {
	switch value := value.(type) {
	case string:
		switch value {
		case Offline, OnlineAvailable, Busy:
			return true
		default:
			return false
		}
	default:
		panic("the accessibility field must be a string")
	}
}

// User represents a user, implements the User interface{}
type User struct {
	Username string `bson:"username" json:"username" valid:"required,alphanum,length(1|30)"`
	Data     `bson:"user_data,inline" valid:"required"`
}

// Data contains the user's data
type Data struct {
	Status               string    `bson:"status" json:"status" valid:"required,validateUserStatus"`
	Email                string    `bson:"email" json:"email" valid:"required,email"`
	Firstname            string    `bson:"firstName" json:"firstName" valid:"required,alpha,length(0|50)"`
	Lastname             string    `bson:"lastName" json:"lastName" valid:"required,alpha,length(0|50)"`
	DefaultAccessibility string    `bson:"defaultAccess" json:"defaultAccess" valid:"required,goalAccessField"`
	PictureURL           *string   `bson:"pictureUrl,omitempty" json:"pictureUrl,omitempty" valid:"optional,url"`
	LastUpdated          time.Time `bson:"lastUpdated" json:"lastUpdated" valid:"-"`
	HasUpdate            bool      `bson:"hasUpdated" json:"hasUpdated" valid:"-"`
}

// CreateGoal creates a new goal
func (c *User) CreateGoal(goalCol *mgo.Collection, g *goal.Goal) (utils.HexID, error) {
	id := utils.HexID{ObjectId: bson.NewObjectId()}
	g.Username, g.HexID = c.Username, id
	err := goalCol.Insert(g)
	if err != nil {
		// Does not catch duplication error
		return id, err
	}
	return id, nil
}

// DeleteGoal deletes a goal
func (c *User) DeleteGoal(goalCol *mgo.Collection, id utils.HexID) error {
	err := goalCol.Remove(bson.M{
		"username": c.Username,
		"_id":      id.ObjectId,
	})
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.NewNotFound("goal", fmt.Sprintf("%s,%s", c.Username, id))
		}
		return err
	}
	return nil
}

// UpdateGoal updates a goal
func (c *User) UpdateGoal(goalCol *mgo.Collection, g *goal.Goal, id utils.HexID) error {
	g.Username, g.HexID = c.Username, id
	err := goalCol.Update(bson.M{
		"username": c.Username,
		"_id":      id,
	}, g)
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.NewNotFound("goal", fmt.Sprintf("%s,%s", c.Username, id))
		}
		return err
	}
	return nil
}

// RetrieveGoal gets the goal
func (c *User) RetrieveGoal(goalCol *mgo.Collection, id utils.HexID) (goal.Goal, error) {
	var g goal.Goal
	err := goalCol.FindId(id.ObjectId).One(&g)
	if err != nil {
		if err == mgo.ErrNotFound {
			return goal.Goal{}, errors.NewNotFound("goal", fmt.Sprintf("%s,%s", c.Username, id))
		}
		return goal.Goal{}, err
	}
	return g, nil
}

// RetriveGoals retrieves all the goals
func (c *User) RetriveGoals(goalCol *mgo.Collection, limit, offset int) ([]goal.Goal, error) {
	var gs []goal.Goal
	q := goalCol.Find(bson.M{
		"username": c.Username,
	})
	if limit > 0 {
		q.Limit(limit)
	}
	if offset > 0 {
		q.Skip(offset)
	}
	err := q.All(&gs)
	if err != nil {
		return nil, err
	}
	return gs, nil
}

// StoredUser encapsulates user's password
type storedUser struct {
	User           `bson:"user,inline" valid:"required"`
	Authentication `bson:"authentication,inline" json:"-" valid:"required"`
}

// Authentication stores authentication data
type Authentication struct {
	Password string `bson:"password" json:"password,omitempty" valid:"hexadecimal,optional"`
	Salt     string `bson:"salt" json:"-" valid:"alphanum,optional"`
}

// Retrieve retrieves the user with the username
func (c *User) Retrieve(userC *mgo.Collection, username string) error {
	err := userC.Find(bson.M{
		"username": username,
	}).One(c)
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.NewNotFound("user", username)
		}
		return err
	}
	return nil
}

// Create creates a new user using the provided password
func Create(c *User, userC *mgo.Collection, password string) error {
	pass, salt, err := randomSaltEncryption(password)
	if err != nil {
		return errors.NewValidate(err.Error())
	}
	usr := storedUser{
		User: *c,
		Authentication: Authentication{
			Password: pass,
			Salt:     salt,
		},
	}
	err = userC.Insert(usr)
	if err != nil {
		if mgo.IsDup(err) {
			return errors.NewDuplicated("user", c.Username)
		}
		return err
	}
	return nil
}

// Update updates a user data
func Update(u *User, userC *mgo.Collection, username string) error {
	u.Username = username
	err := userC.Update(bson.M{
		"username": username,
	}, bson.M{
		"$set": u,
	})
	if err != nil {
		if err == mgo.ErrNotFound {
			return errors.NewNotFound("user", username)
		}
		return err
	}
	return nil
}

// Delete deletes an user
func Delete(userC *mgo.Collection, username, password string) error {
	ok, err := Authenticate(userC, username, password)
	if err != nil {
		return err
	}
	if ok {
		err = userC.Remove(bson.M{
			"username": username,
		})
		if err != nil {
			if err == mgo.ErrNotFound {
				return errors.NewNotFound("user", username)
			}
			return err
		}
		return nil
	}
	return errors.ErrAuthentication
}

// Authenticate authenticates the user
func Authenticate(userCol *mgo.Collection, username, password string) (bool, error) {
	var auth Authentication
	err := userCol.Find(bson.M{
		"username": username,
	}).Select(bson.M{
		"password": 1,
		"salt":     1,
	}).One(&auth)
	if err != nil {
		if err == mgo.ErrNotFound {
			return false, errors.NewNotFound("user", username)
		}
		return false, err
	}
	if auth.Password != encryptPassword(auth.Salt, password) {
		return false, nil
	}
	return true, nil
}

// ChangePassword changes the user's password
func ChangePassword(userC *mgo.Collection, username, old, password string) error {
	ok, err := Authenticate(userC, username, old)
	if err != nil {
		return err
	}
	if ok {
		salt, pass, err := randomSaltEncryption(password)
		if err != nil {
			return errors.NewValidate(err.Error())
		}
		err = userC.Update(bson.M{
			"username": username,
		}, bson.M{
			"$set": bson.M{
				"password": pass,
				"salt":     salt,
			},
		})
		if err != nil {
			if err == mgo.ErrNotFound {
				return errors.NewNotFound("user", username)
			}
			return err
		}
		return nil
	}
	return errors.ErrAuthentication
}
