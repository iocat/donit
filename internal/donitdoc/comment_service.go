package donitdoc

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/comments"
	"github.com/iocat/donit/internal/donitdoc/utils"
	"gopkg.in/mgo.v2"
)

type commentService service

// CommentService represents a service that deals with the comment data model
type CommentService interface {
	Create(context.Context, *comments.Comment) (string, error)
	Read(context.Context, string) (*comments.Comment, error)
	Update(context.Context, *comments.Comment) error
	Delete(context.Context, string) error
}

// NewComment creates a CommentService
func NewComment(l log.Logger, db *mgo.Database) (CommentService, error) {
	return &commentService{
		l:  l,
		db: db,
	}, nil
}

func (s *commentService) Create(ctx context.Context, g *comments.Comment) (string, error) {
	id, err := comments.CreateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *commentService) Read(ctx context.Context, id string) (*comments.Comment, error) {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.ReadComment", "result", err)
		return nil, err
	}
	return comments.ReadDoc(s.l, lid, s.db, oid)
}

func (s *commentService) Update(ctx context.Context, g *comments.Comment) error {
	return comments.UpdateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
}
func (s *commentService) Delete(ctx context.Context, id string) error {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.DeleteComment", "result", err)
		return err
	}
	return comments.DeleteDoc(s.l, lid, s.db, oid)
}
