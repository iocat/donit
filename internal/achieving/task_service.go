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

	"github.com/iocat/donit/internal/donitdoc/tasks"
	"gopkg.in/mgo.v2"
)

// TaskService represents a service that deals with the task data model
type TaskService interface {
	Create(context.Context, *tasks.Task) (string, error)
	Read(context.Context, string) (*tasks.Task, error)
	Update(context.Context, *tasks.Task) error
	Delete(context.Context, string) error
}

type taskService service

// NewTask creates a TaskService
func NewTask(l log.Logger, db *mgo.Database) (TaskService, error) {
	return &taskService{
		l:  l,
		db: db,
	}, nil
}

func (s *taskService) Create(ctx context.Context, g *tasks.Task) (string, error) {
	id, err := tasks.CreateDoc(s.l, logID(ctx), s.db, g)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *taskService) Read(ctx context.Context, id string) (*tasks.Task, error) {
	lid := logID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.ReadTask", "result", err)
		return nil, err
	}
	return tasks.ReadDoc(s.l, lid, s.db, oid)
}

func (s *taskService) Update(ctx context.Context, g *tasks.Task) error {
	return tasks.UpdateDoc(s.l, logID(ctx), s.db, g)
}
func (s *taskService) Delete(ctx context.Context, id string) error {
	lid := logID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.DeleteTask", "result", err)
		return err
	}
	return tasks.DeleteDoc(s.l, lid, s.db, oid)
}
