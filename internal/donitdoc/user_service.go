package donitdoc

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/users"
	"github.com/iocat/donit/internal/donitdoc/utils"
	"gopkg.in/mgo.v2"
)

type userService service

// UserService represents a service that deals with the user data model
type UserService interface {
	Create(context.Context, *users.User, string) error
	Read(context.Context, string) (*users.User, error)
	Update(context.Context, *users.User) error
	Delete(context.Context, string) error

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
	userCollection, err := utils.GetMGOCollectionName(utils.User)
	if err != nil {
		return fmt.Errorf("get collection name for user: %s", err)
	}
	s.db.C(userCollection).EnsureIndex(mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: false,
		Sparse:     true,
	})
	return nil
}

func (s *userService) Create(ctx context.Context, user *users.User, pass string) error {
	return users.CreateDoc(s.l, utils.MustGetLogID(ctx), s.db, user, pass)
}

func (s *userService) Read(ctx context.Context, username string) (*users.User, error) {
	return users.ReadDoc(s.l, utils.MustGetLogID(ctx), s.db, username)
}

func (s *userService) Update(ctx context.Context, user *users.User) error {
	return users.UpdateDoc(s.l, utils.MustGetLogID(ctx), s.db, user)
}

func (s *userService) Delete(ctx context.Context, username string) error {
	return users.DeleteDoc(s.l, utils.MustGetLogID(ctx), s.db, username)
}

func (s *userService) Authenticate(ctx context.Context, username, password string) (bool, error) {
	return users.ValidatePassword(s.l, utils.MustGetLogID(ctx), s.db, username, password)
}

func (s *userService) ChangePassword(ctx context.Context, username, oldpass, newpass string) error {
	return users.ChangePassword(s.l, utils.MustGetLogID(ctx), s.db, username, oldpass, newpass)
}
