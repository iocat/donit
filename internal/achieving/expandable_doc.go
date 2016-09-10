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

// ExpandableUser represents a User object that contains every related data
// Expandable forms a tree structure inside the struct to support efficient
// retrieval on the backend
type ExpandableUser struct {
	User
	Goals []ExpandableGoal `json:"goals"`
}

// ExpandableGoal represents a Expandable goal
type ExpandableGoal struct {
	Goal
	Tasks []Achievable `json:"tasks"`
}

// NewExpandableUser creates an ExpandableUser
func NewExpandableUser(user User) *ExpandableUser {
	return &ExpandableUser{
		User: user,
	}
}

// Retrieve gets all the goals (a GOALDEN RETRIEVER) and expand the goal if needed
func (e *ExpandableUser) Retrieve() error {
	goals, err := e.RetrieveGoals(0, 0)
	if err != nil {
		return err
	}
	e.Goals = make([]ExpandableGoal, len(goals))
	for _, goal := range goals {
		eg := ExpandableGoal{
			Goal: goal,
		}
		e.Goals = append(e.Goals, eg)
		if err := eg.Retrieve(); err != nil {
			return err
		}
	}
	return nil
}

// Retrieve gets all the tasks associated with this goal and store them
func (e *ExpandableGoal) Retrieve() error {
	var err error
	e.Tasks, err = e.RetrieveAchievable(0, 0)
	if err != nil {
		return err
	}
	return nil
}

// NewExpandableGoal creates an ExpandableGoal
func NewExpandableGoal(g Goal) *ExpandableGoal {
	return &ExpandableGoal{
		Goal: g,
	}
}
