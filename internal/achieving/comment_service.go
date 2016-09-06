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
	id, err := comments.CreateDoc(s.l, logID(ctx), s.db, g)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *commentService) Read(ctx context.Context, id string) (*comments.Comment, error) {
	lid := logID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.ReadComment", "result", err)
		return nil, err
	}
	return comments.ReadDoc(s.l, lid, s.db, oid)
}

func (s *commentService) Update(ctx context.Context, g *comments.Comment) error {
	return comments.UpdateDoc(s.l, logID(ctx), s.db, g)
}
func (s *commentService) Delete(ctx context.Context, id string) error {
	lid := logID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.DeleteComment", "result", err)
		return err
	}
	return comments.DeleteDoc(s.l, lid, s.db, oid)
}
