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

// Package docgroup is an abstraction layer above donitdoc and other
// service oriented packages. Its main purpose is to group document in a tree
// hierarchy and to support efficient retrieval of documents
package docgroup

import (
	"context"

	"github.com/iocat/donit/internal/achiever"
	"github.com/iocat/donit/internal/achieving"
	"github.com/iocat/donit/internal/donitdoc/comments"
	"github.com/iocat/donit/internal/donitdoc/followers"
	"github.com/iocat/donit/internal/donitdoc/goals"
	"github.com/iocat/donit/internal/donitdoc/habits"
	"github.com/iocat/donit/internal/donitdoc/tasks"
	"github.com/iocat/donit/internal/donitdoc/users"
)

// ExpandedReader represents a reader that can expand the entities down the
// tree hierarchy of the document
type ExpandedReader interface {
	Reader
	ExpandRead(context.Context) error
}

// Reader represents a document reader
type Reader interface {
	Read(context.Context) error
}

// Habit represents a habit
type Habit struct {
	*habits.Habit
	achieving.HabitService
}

// Task represents a task
type Task struct {
	*tasks.Task
	achieving.TaskService
}

// Comment represents a comment
type Comment struct {
	*comments.Comment
	achieving.CommentService
}

// Follower represents a follower
type Follower struct {
	*followers.Follower
	achiever.FollowerService
}

// User represents a user
type User struct {
	*users.User

	Followers []Follower `json:"followers"`
	Goals     []Goal     `json:"goals"`

	achiever.UserService

	achiever.FollowerService
	achieving.GoalService

	achieving.TaskService
	achieving.HabitService
	achieving.CommentService
}

// Goal represents a goal
type Goal struct {
	*goals.Goal

	Tasks    []Task    `json:"tasks"`
	Habits   []Habit   `json:"habits"`
	Comments []Comment `json:"comments"`

	achieving.GoalService

	achieving.TaskService
	achieving.HabitService
	achieving.CommentService
}

// ExpandRead implements ExpandedReader
func (user *User) ExpandRead(ctx context.Context) error {
	return user.read(ctx, true)
}

// Read implements Reader
func (user *User) Read(ctx context.Context) error {
	return user.read(ctx, false)
}

// Read implements Reader
func (g *Goal) Read(ctx context.Context) error {
	return g.read(ctx, false)
}

// ExpandRead implements ExpandedReader
func (g *Goal) ExpandRead(ctx context.Context) error {
	return g.read(ctx, true)
}

func (h *Habit) Read(ctx context.Context) error {
	return h.read(ctx)
}

func (t *Task) Read(ctx context.Context) error {
	return t.read(ctx)
}

func (c *Comment) Read(ctx context.Context) error {
	return c.read(ctx)
}

func (user *User) read(ctx context.Context, expand bool) error {
	var err error
	if user.User, err = user.UserService.Read(ctx, user.Username); err != nil {
		return err
	}
	if !expand {
		return nil
	}
	// Expand the follower list
	followers, err := user.AllFollowers(ctx, user.Username, -1, -1)
	if err != nil {
		return err
	}
	user.Followers = make([]Follower, len(followers))
	for i := range followers {
		user.Followers[i] = Follower{Follower: &followers[i]}
	}
	// Expand the goal list
	goals, err := user.AllGoals(ctx, user.Username, -1, -1)
	if err != nil {
		return err
	}
	user.Goals = make([]Goal, len(goals))
	for i := range goals {
		// copy the actual data
		user.Goals[i].Goal = &goals[i]
		// Copy the service
		user.Goals[i].GoalService = user.GoalService
		user.Goals[i].TaskService = user.TaskService
		user.Goals[i].HabitService = user.HabitService
		// expandable read the underlying goals
		if err := user.Goals[i].read(ctx, expand); err != nil {
			return err
		}
	}
	return nil
}

func (g *Goal) read(ctx context.Context, expand bool) error {
	var err error
	g.Goal, err = g.GoalService.Read(ctx, g.Goal.ObjectId.Hex())
	if err != nil {
		return err
	}
	if !expand {
		return nil
	}
	// Expand habit
	habits, err := g.AllHabits(ctx, g.Goal.ObjectId.Hex(), -1, -1)
	if err != nil {
		return err
	}
	g.Habits = make([]Habit, len(habits))
	for i := range habits {
		g.Habits[i] = Habit{
			Habit:        &habits[i],
			HabitService: g.HabitService,
		}
	}
	// Expand task
	tasks, err := g.AllTasks(ctx, g.Goal.ObjectId.Hex(), -1, -1)
	if err != nil {
		return err
	}
	g.Tasks = make([]Task, len(habits))
	for i := range tasks {
		g.Tasks[i] = Task{
			Task:        &tasks[i],
			TaskService: g.TaskService,
		}
	}
	// Expand comments
	comments, err := g.AllComments(ctx, g.Goal.ObjectId.Hex(), -1, -1)
	if err != nil {
		return err
	}
	g.Comments = make([]Comment, len(comments))
	for i := range comments {
		g.Comments[i] = Comment{
			Comment:        &comments[i],
			CommentService: g.CommentService,
		}
	}
	return nil
}

func (h *Habit) read(ctx context.Context) error {
	var err error
	h.Habit, err = h.HabitService.Read(ctx, h.Habit.ObjectId.Hex())
	if err != nil {
		return err
	}
	return nil
}

func (t *Task) read(ctx context.Context) error {
	var err error
	t.Task, err = t.TaskService.Read(ctx, t.Task.ObjectId.Hex())
	if err != nil {
		return err
	}
	return nil
}

func (c *Comment) read(ctx context.Context) error {
	var err error
	c.Comment, err = c.CommentService.Read(ctx, c.Comment.ObjectId.Hex())
	if err != nil {
		return err
	}
	return nil
}
