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

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/comments"
	"github.com/iocat/donit/internal/donitdoc/goals"
	"github.com/iocat/donit/internal/donitdoc/habits"
	"github.com/iocat/donit/internal/donitdoc/tasks"
	"gopkg.in/mgo.v2"
)

// GoalService represents a service that deals with the goal data model
// TODO: implement uncomment methods
type GoalService interface {
	Create(context.Context, *goals.Goal) (string, error)
	Read(context.Context, string) (*goals.Goal, error)
	Update(context.Context, *goals.Goal) error
	Delete(context.Context, string) error

	AllComments(context.Context, string, int, int) ([]comments.Comment, error)
	AllTasks(context.Context, string, int, int) ([]tasks.Task, error)
	AllHabits(context.Context, string, int, int) ([]habits.Habit, error)
}

type goalService service

// NewGoal creates a GoalService
func NewGoal(l log.Logger, db *mgo.Database) (GoalService, error) {
	return &goalService{
		l:  l,
		db: db,
	}, nil
}

func (s *goalService) Create(ctx context.Context, g *goals.Goal) (string, error) {
	id, err := goals.CreateDoc(s.l, logID(ctx), s.db, g)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *goalService) Read(ctx context.Context, id string) (*goals.Goal, error) {
	lid := logID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.ReadGoal", "result", err)
		return nil, err
	}
	return goals.ReadDoc(s.l, lid, s.db, oid)
}

func (s *goalService) Update(ctx context.Context, g *goals.Goal) error {
	return goals.UpdateDoc(s.l, logID(ctx), s.db, g)
}

func (s *goalService) Delete(ctx context.Context, id string) error {
	lid := logID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.DeleteGoal", "result", err)
		return err
	}
	if err := comments.DeleteAllDocsOfGoal(s.l, lid, s.db, oid); err != nil {
		return err
	}
	if err := tasks.DeleteAllDocsOfGoal(s.l, lid, s.db, oid); err != nil {
		return err
	}
	if err := habits.DeleteAllDocsOfGoal(s.l, lid, s.db, oid); err != nil {
		return err
	}
	if err := goals.DeleteDoc(s.l, lid, s.db, oid); err != nil {
		return err
	}
	return nil
}

func (s *goalService) AllComments(ctx context.Context, goal string, lim, off int) ([]comments.Comment, error) {
	lid := logID(ctx)
	oid, err := getID(goal)
	if err != nil {
		s.l.Log("ctx", lid, "op", "Goalservice.AllComments", "result", err)
		return nil, err
	}
	return comments.ALlDocsOfGoal(s.l, lid, s.db, oid, lim, off)
}
func (s *goalService) AllTasks(ctx context.Context, goal string, lim, off int) ([]tasks.Task, error) {
	lid := logID(ctx)
	oid, err := getID(goal)
	if err != nil {
		s.l.Log("ctx", lid, "op", "Goalservice.AllTasks", "result", err)
		return nil, err
	}
	return tasks.AllDocsOfGoal(s.l, lid, s.db, oid, lim, off)
}
func (s *goalService) AllHabits(ctx context.Context, goal string, lim int, off int) ([]habits.Habit, error) {
	lid := logID(ctx)
	oid, err := getID(goal)
	if err != nil {
		s.l.Log("ctx", lid, "op", "Goalservice.AllHabits", "result", err)
		return nil, err
	}
	return habits.AllDocsOfGoal(s.l, lid, s.db, oid, lim, off)
}
