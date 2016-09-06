package achiever

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/errors"
	"github.com/iocat/donit/internal/donitdoc/followers"
	"github.com/iocat/donit/internal/donitdoc/goals"
	"github.com/iocat/donit/internal/donitdoc/users"
	"github.com/iocat/donit/internal/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type userService service

type service struct {
	l  log.Logger
	db *mgo.Database
}

func getID(id string) (bson.ObjectId, error) {
	if !bson.IsObjectIdHex(id) {
		return "", errors.NewValidate(fmt.Sprintf("object id %s is invalid ", id))
	}
	return bson.ObjectIdHex(id), nil
}

func init() {
	logID = utils.MustGetLogID
}

var logID func(context.Context) string

// UserService represents a service that deals with the user data model
type UserService interface {
	Create(context.Context, *users.User, string) error
	Read(context.Context, string) (*users.User, error)
	Update(context.Context, *users.User) error
	Delete(context.Context, string) error

	AllFollowers(context.Context, string, int, int) ([]followers.Follower, error)
	AllGoals(context.Context, string, int, int) ([]goals.Goal, error)

	Authenticate(context.Context, string, string) (bool, error)
	ChangePassword(context.Context, string, string, string) error
}

// NewUser creates an UserService
func NewUser(l log.Logger, db *mgo.Database) (UserService, error) {
	us := userService{
		l:  l,
		db: db,
	}
	if err := us.dbEnsureUserIndex(); err != nil {
		return nil, err
	}
	return &us, nil
}

func (s *userService) dbEnsureUserIndex() error {
	s.db.C(users.CollectionName).EnsureIndex(mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: false,
		Sparse:     true,
	})
	return nil
}

func (s *userService) AllGoals(ctx context.Context, user string, limit, offset int) ([]goals.Goal, error) {
	return goals.AllDocsOfUser(s.l, logID(ctx), s.db, user, limit, offset)
}

func (s *userService) Create(ctx context.Context, user *users.User, pass string) error {
	return users.CreateDoc(s.l, logID(ctx), s.db, user, pass)
}

func (s *userService) Read(ctx context.Context, username string) (*users.User, error) {
	return users.ReadDoc(s.l, logID(ctx), s.db, username)
}

func (s *userService) Update(ctx context.Context, user *users.User) error {
	return users.UpdateDoc(s.l, logID(ctx), s.db, user)
}

func (s *userService) Delete(ctx context.Context, username string) error {
	if err := followers.DeleteAllDocsOfUser(s.l, logID(ctx), s.db, username); err != nil {
		return err
	}
	if err := goals.DeleteAllDocsOfUser(s.l, logID(ctx), s.db, username); err != nil {
		return err
	}
	if err := users.DeleteDoc(s.l, logID(ctx), s.db, username); err != nil {
		return err
	}
	return nil
}

func (s *userService) AllFollowers(ctx context.Context, username string, limit, offset int) ([]followers.Follower, error) {
	return followers.AllDocsOfUser(s.l, logID(ctx), s.db, username, limit, offset)
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (bool, error) {
	return users.ValidatePassword(s.l, logID(ctx), s.db, username, password)
}

func (s *userService) ChangePassword(ctx context.Context, username, oldpass, newpass string) error {
	return users.ChangePassword(s.l, logID(ctx), s.db, username, oldpass, newpass)
}
