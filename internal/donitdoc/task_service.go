package donitdoc

import (
	"context"

	"github.com/go-kit/kit/log"

	"github.com/iocat/donit/internal/donitdoc/tasks"
	"github.com/iocat/donit/internal/donitdoc/utils"
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
	id, err := tasks.CreateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *taskService) Read(ctx context.Context, id string) (*tasks.Task, error) {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.ReadTask", "result", err)
		return nil, err
	}
	return tasks.ReadDoc(s.l, lid, s.db, oid)
}

func (s *taskService) Update(ctx context.Context, g *tasks.Task) error {
	return tasks.UpdateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
}
func (s *taskService) Delete(ctx context.Context, id string) error {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.DeleteTask", "result", err)
		return err
	}
	return tasks.DeleteDoc(s.l, lid, s.db, oid)
}
