package donitdoc

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/goals"
	"github.com/iocat/donit/internal/donitdoc/utils"
	"gopkg.in/mgo.v2"
)

// GoalService represents a service that deals with the goal data model
type GoalService interface {
	Create(context.Context, *goals.Goal) (string, error)
	Read(context.Context, string) (*goals.Goal, error)
	Update(context.Context, *goals.Goal) error
	Delete(context.Context, string) error
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
	id, err := goals.CreateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *goalService) Read(ctx context.Context, id string) (*goals.Goal, error) {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.ReadGoal", "result", err)
		return nil, err
	}
	return goals.ReadDoc(s.l, lid, s.db, oid)
}

func (s *goalService) Update(ctx context.Context, g *goals.Goal) error {
	return goals.UpdateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
}
func (s *goalService) Delete(ctx context.Context, id string) error {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.DeleteGoal", "result", err)
		return err
	}
	return goals.DeleteDoc(s.l, lid, s.db, oid)
}
