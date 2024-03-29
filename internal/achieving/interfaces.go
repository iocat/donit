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

package achieving

// Achievable represents an achieveable task
type Achievable interface {
	HasAchieved() bool
	IsRepetitive() bool
}

// Goal represents a goal that has the goal data
// Goal also acts as a Achievable object container
type Goal interface {
	// AddAchievableTask adds the task
	AddAchievable(Achievable) (string, error)
	// RemoveAchievableTask removes the task
	RemoveAchievable(string) error
	// UpdateAchievableTask updates the task
	UpdateAchievable(Achievable, string) error

	// RetriveAchievableTask gets a list of achievable task
	RetrieveAchievables(limit, offset int) ([]Achievable, error)
}

// User represents am user object, which should be containing the user data
// User is also a Goal container that has goals' CRUD operations
type User interface {
	// Create goal creates a new goal
	CreateGoal(Goal) (string, error)
	// DeleteGoal deletes a goal
	DeleteGoal(string) error
	// UpdateGoal updates a goal
	UpdateGoal(Goal, string) error
	// RetrieveGoal retrieves a goal
	RetrieveGoal(string) (Goal, error)

	// RetrieveGoals get all the goal from this user
	RetrieveGoals(limit, offset int) ([]Goal, error)
}

// UserStore represents a storage of user, it does not contain the user data
// UserStore allows operations on UserStore
type UserStore interface {
	// RetriveUser retrives the user and expand the user data as needed
	RetrieveUser(string) (User, error)
	// CreateNewUser creates a new user using the provided username and password
	CreateNewUser(User, string) (string, error)
	// DeleteUser deletes a user using the provided username and password
	DeleteUser(string, string) error
	// UpdateUser updates the user information
	UpdateUser(User, string) error

	// Authenticate authenticates the username and password
	Authenticate(string, string) (bool, error)
	// ChangePassword changes a user password
	ChangePassword(string, string, string) error
}
