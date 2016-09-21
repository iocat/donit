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

package concreteachieving

import (
	"fmt"

	"github.com/iocat/donit/internal/achieving"
	"github.com/iocat/donit/internal/achieving/internal/user"
	"gopkg.in/mgo.v2"
)

// Store implements the achieving.UserStore
type Store struct {
	userCollection       *mgo.Collection `valid:"-"`
	goalCollection       *mgo.Collection `valid:"-"`
	achievableCollection *mgo.Collection `valid:"-"`
}

// NewStore creates a new UserStore
func NewStore(user, goal, task *mgo.Collection) *Store {
	return &Store{
		userCollection:       user,
		goalCollection:       goal,
		achievableCollection: task,
	}
}

// RetrieveUser implements userStore
func (s Store) RetrieveUser(username string) (achieving.User, error) {
	u := user.User{}
	err := u.Retrieve(s.userCollection, username)
	if err != nil {
		return nil, err
	}
	cu := User{
		User:                 u,
		goalCollection:       s.goalCollection,
		achievableCollection: s.achievableCollection,
	}
	return &cu, nil
}

// CreateNewUser creates a new user using the provided username and password
func (s Store) CreateNewUser(u achieving.User, password string) error {
	if u, ok := u.(*User); ok {
		err := user.Create(&(u.User), s.userCollection, password)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("wrong data type, expect *concreteachieving.User, got %T", u)
}

// DeleteUser deletes a user using the provided username and password
func (s Store) DeleteUser(username, password string) error {
	return user.Delete(s.userCollection, username, password)
}

// Authenticate authenticates the username and password
func (s Store) Authenticate(username, password string) (bool, error) {
	return user.Authenticate(s.userCollection, username, password)
}

// ChangePassword changes a user password
func (s Store) ChangePassword(username, oldpass, password string) error {
	return user.ChangePassword(s.userCollection, username, oldpass, password)
}

// UpdateUser updates a user data
func (s Store) UpdateUser(u achieving.User, username string) error {
	if u, ok := u.(*User); !ok {
		err := user.Update(&(u.User), s.userCollection, username)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("wrong data type, expect *concreteachieving.User, got %T", u)
}
