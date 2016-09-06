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

package utils

import "gopkg.in/mgo.v2"

const (
	// User is the collection index for User
	User Collection = iota
	// Goal is the collection index for Goal
	Goal
	// Habit is the collection index for Habit
	Habit
	// Task is the collection index for Task
	Task
	// Comment is the collection index for Comment
	Comment
	// Follower is the colelction index for Follower
	Follower
)

// Collection serializes collections
type Collection byte

// Collections contains a list of collections
var Collections = []Collection{User, Goal, Habit, Task, Comment}

// MGOName returns the MGO collection's name corresponding to the collection code
func (c Collection) MGOName() string {
	switch c {
	case User:
		return "users"
	case Goal:
		return "goals"
	case Habit:
		return "habits"
	case Task:
		return "tasks"
	case Comment:
		return "comments"
	case Follower:
		return "followers"
	default:
		return ""
	}
}

// Collection creates/returns a collection from the database
func (c Collection) Collection(db *mgo.Database) *mgo.Collection {
	return db.C(c.MGOName())
}

// Query decorates mgo.Query
type Query struct{ *mgo.Query }

// Skip skips a certain number of document. If offset < 0, the query
// doesn't skip document.
func (q *Query) Skip(offset int) *Query {
	if offset < 0 {
		return q
	}
	return q.Skip(offset)
}

// Limit limits the number of returned documents. If limit <= 0, the
// query doesn't limit the number of document
func (q *Query) Limit(limit int) *Query {
	if limit <= 0 {
		return q
	}
	return q.Limit(limit)
}
