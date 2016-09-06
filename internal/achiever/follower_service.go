package achiever

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/iocat/donit/internal/donitdoc/followers"
	"gopkg.in/mgo.v2"
)

type followerService service

// FollowerService represents a follower service
type FollowerService interface {
	Follows(context.Context, *followers.Follower) error
	Unfollows(context.Context, string, string) error
}

// Follows let a follower follows the username
func (f *followerService) Follows(ctx context.Context, fl *followers.Follower) error {
	return followers.CreateDoc(f.l, logID(ctx), f.db, fl)
}

// Unfollows lets the follower unfollows the username
func (f *followerService) Unfollows(ctx context.Context, username, follower string) error {
	return followers.DeleteDoc(f.l, logID(ctx), f.db, username, follower)
}

func (f *followerService) dbEnsureFollowerIndex() error {
	f.db.C(followers.CollectionName).EnsureIndex(mgo.Index{
		Key:        []string{"username", "follower"},
		Unique:     true,
		DropDups:   true,
		Background: false,
		Sparse:     true,
	})
	return nil
}

// NewFollower creates a new follower service
func NewFollower(l log.Logger, db *mgo.Database) (FollowerService, error) {
	fs := &followerService{
		l:  l,
		db: db,
	}
	if err := fs.dbEnsureFollowerIndex(); err != nil {
		return nil, err
	}
	return fs, nil
}
