package donitdoc

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/comments"
	"github.com/iocat/donit/internal/donitdoc/errors"
	"github.com/iocat/donit/internal/donitdoc/goals"
	"github.com/iocat/donit/internal/donitdoc/habits"
	"github.com/iocat/donit/internal/donitdoc/tasks"
	"github.com/iocat/donit/internal/donitdoc/users"
	"github.com/iocat/donit/internal/donitdoc/utils"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// DocumentService is the facade API for donit's document store
// NOTE: This is a fat interface for exported API. This interface is not expected
// to be changed
type DocumentService interface {
	CreateUser(context.Context, *users.User, string) error
	ReadUser(context.Context, string) (*users.User, error)
	UpdateUser(context.Context, *users.User) error
	DeleteUser(context.Context, string) error

	AuthenticateUser(context.Context, string, string) (bool, error)
	ChangeUserPassword(context.Context, string, string, string) error

	CreateGoal(context.Context, *goals.Goal) (string, error)
	ReadGoal(context.Context, string) (*goals.Goal, error)
	UpdateGoal(context.Context, *goals.Goal) error
	DeleteGoal(context.Context, string) error

	CreateHabit(context.Context, *habits.Habit) (string, error)
	ReadHabit(context.Context, string) (*habits.Habit, error)
	UpdateHabit(context.Context, *habits.Habit) error
	DeleteHabit(context.Context, string) error

	CreateTask(context.Context, *tasks.Task) (string, error)
	ReadTask(context.Context, string) (*tasks.Task, error)
	UpdateTask(context.Context, *tasks.Task) error
	DeleteTask(context.Context, string) error

	CreateComment(context.Context, *comments.Comment) (string, error)
	ReadComment(context.Context, string) (*comments.Comment, error)
	UpdateComment(context.Context, *comments.Comment) error
	DeleteComment(context.Context, string) error
}

// service implements Service interface
type service struct {
	dbSession *mgo.Session
	db        *mgo.Database
	l         log.Logger
}

var dbDialTimeout = 2 * time.Second

// New creates a new DocumentService
func New(l log.Logger, DBUrl string, DBName string) (DocumentService, error) {
	db, err := mgo.DialWithTimeout(DBUrl, dbDialTimeout)
	if err != nil {
		return nil, fmt.Errorf("dial to mongodb: %s", err)
	}
	s := service{
		dbSession: db,
		l:         l,
	}
	s.db = db.DB(DBName)
	if err := s.dbSetup(); err != nil {
		return nil, fmt.Errorf("mongodb set up: %s", err)
	}
	return &s, nil
}

func (s *service) dbSetup() error {
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

func (s *service) CreateUser(ctx context.Context, user *users.User, pass string) error {
	return users.CreateDoc(s.l, utils.MustGetLogID(ctx), s.db, user, pass)
}

func (s *service) ReadUser(ctx context.Context, username string) (*users.User, error) {
	return users.ReadDoc(s.l, utils.MustGetLogID(ctx), s.db, username)
}

func (s *service) UpdateUser(ctx context.Context, user *users.User) error {
	return users.UpdateDoc(s.l, utils.MustGetLogID(ctx), s.db, user)
}

func (s *service) DeleteUser(ctx context.Context, username string) error {
	return users.DeleteDoc(s.l, utils.MustGetLogID(ctx), s.db, username)
}

func (s *service) AuthenticateUser(ctx context.Context, username, password string) (bool, error) {
	return users.ValidatePassword(s.l, utils.MustGetLogID(ctx), s.db, username, password)
}

func (s *service) ChangeUserPassword(ctx context.Context, username, oldpass, newpass string) error {
	return users.ChangePassword(s.l, utils.MustGetLogID(ctx), s.db, username, oldpass, newpass)
}

func getID(id string) (bson.ObjectId, error) {
	if !bson.IsObjectIdHex(id) {
		return "", errors.NewValidate(fmt.Sprintf("object id %s is invalid ", id))
	}
	return bson.ObjectIdHex(id), nil
}

func (s *service) CreateGoal(ctx context.Context, g *goals.Goal) (string, error) {
	id, err := goals.CreateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *service) ReadGoal(ctx context.Context, id string) (*goals.Goal, error) {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.ReadGoal", "result", err)
		return nil, err
	}
	return goals.ReadDoc(s.l, lid, s.db, oid)
}

func (s *service) UpdateGoal(ctx context.Context, g *goals.Goal) error {
	return goals.UpdateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
}
func (s *service) DeleteGoal(ctx context.Context, id string) error {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.DeleteGoal", "result", err)
		return err
	}
	return goals.DeleteDoc(s.l, lid, s.db, oid)
}

func (s *service) CreateTask(ctx context.Context, g *tasks.Task) (string, error) {
	id, err := tasks.CreateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *service) ReadTask(ctx context.Context, id string) (*tasks.Task, error) {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.ReadTask", "result", err)
		return nil, err
	}
	return tasks.ReadDoc(s.l, lid, s.db, oid)
}

func (s *service) UpdateTask(ctx context.Context, g *tasks.Task) error {
	return tasks.UpdateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
}
func (s *service) DeleteTask(ctx context.Context, id string) error {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.DeleteTask", "result", err)
		return err
	}
	return tasks.DeleteDoc(s.l, lid, s.db, oid)
}

func (s *service) CreateHabit(ctx context.Context, g *habits.Habit) (string, error) {
	id, err := habits.CreateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *service) ReadHabit(ctx context.Context, id string) (*habits.Habit, error) {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.ReadHabit", "result", err)
		return nil, err
	}
	return habits.ReadDoc(s.l, lid, s.db, oid)
}

func (s *service) UpdateHabit(ctx context.Context, g *habits.Habit) error {
	return habits.UpdateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
}

func (s *service) DeleteHabit(ctx context.Context, id string) error {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.DeleteHabit", "result", err)
		return err
	}
	return habits.DeleteDoc(s.l, lid, s.db, oid)
}

func (s *service) CreateComment(ctx context.Context, g *comments.Comment) (string, error) {
	id, err := comments.CreateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
	if err != nil {
		return "", err
	}
	return id.Hex(), nil
}

func (s *service) ReadComment(ctx context.Context, id string) (*comments.Comment, error) {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.ReadComment", "result", err)
		return nil, err
	}
	return comments.ReadDoc(s.l, lid, s.db, oid)
}

func (s *service) UpdateComment(ctx context.Context, g *comments.Comment) error {
	return comments.UpdateDoc(s.l, utils.MustGetLogID(ctx), s.db, g)
}
func (s *service) DeleteComment(ctx context.Context, id string) error {
	lid := utils.MustGetLogID(ctx)
	oid, err := getID(id)
	if err != nil {
		s.l.Log("ctx", lid, "op", "DocumentService.DeleteComment", "result", err)
		return err
	}
	return comments.DeleteDoc(s.l, lid, s.db, oid)
}
