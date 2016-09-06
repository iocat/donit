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
	"github.com/iocat/donit/internal/donitdoc/habits"
	"gopkg.in/mgo.v2"
)

// HabitService represents a service that deals with the habit data model
type HabitService interface {
	Create(context.Context, *habits.Habit) (string, error)
	Read(context.Context, string) (*habits.Habit, error)
	Update(context.Context, *habits.Habit) error
	Delete(context.Context, string) error
}

type habitService service

// NewHabit creates a HabitService
func NewHabit(l log.Logger, db *mgo.Database) (HabitService, error) {
	return &habitService{
		l:  l,
		db: db,
	}, nil
}

func (s *habitService) Create(ctx context.Context, g *habits.Habit) (string, error) {
	id, err := habits.CreateDoc(s.l, logID(ctx), s.db, g)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *habitService) Read(ctx context.Context, id string) (*habits.Habit, error) {
	lid := logID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.ReadHabit", "result", err)
		return nil, err
	}
	return habits.ReadDoc(s.l, lid, s.db, oid)
}

func (s *habitService) Update(ctx context.Context, g *habits.Habit) error {
	return habits.UpdateDoc(s.l, logID(ctx), s.db, g)
}

func (s *habitService) Delete(ctx context.Context, id string) error {
	lid := logID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.DeleteHabit", "result", err)
		return err
	}
	return habits.DeleteDoc(s.l, lid, s.db, oid)
}
