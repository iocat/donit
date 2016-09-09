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
	"github.com/iocat/donit/internal/achieving"
	"github.com/iocat/donit/internal/achieving/internal/user"
	"gopkg.in/mgo.v2"
)

// Store implements the achieving.UserStore
type Store struct {
	userCollection       *mgo.Collection
	goalCollection       *mgo.Collection
	achievableCollection *mgo.Collection
}

// RetrieveUser implements userStore
func (us Store) RetrieveUser(username string, expand bool) (achieving.User, error) {
	u := user.User{}
	err := u.Retrieve(us.userCollection, username)
	if err != nil {
		return nil, err
	}
	cu := User{
		User:                 u,
		goalCollection:       us.goalCollection,
		achievableCollection: us.achievableCollection,
	}
	return &cu, nil
}

// CreateNewUser creates a new user using the provided username and password
func (us Store) CreateNewUser(user achieving.User, password string) error {

	return nil
}

// DeleteUser deletes a user using the provided username and password
func (us Store) DeleteUser(username, password string) error {

	return nil
}

// Authenticate authenticates the username and password
func (us Store) Authenticate(username, password string) (bool, error) {

	return false, nil
}

// ChangePassword changes a user password
func (us Store) ChangePassword(username, password string) error {
	return nil
}
