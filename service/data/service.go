package data

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// service implements the Service interface
type service struct {
	*database
}

func (s service) Item(cname string) (Item, error) {
	switch cname {
	case CollectionUser:
		return &User{}, nil
	case CollectionGoal:
		return &Goal{}, nil
	case CollectionFollower:
		return &Follower{}, nil
	case CollectionHabit:
		return &Habit{}, nil
	case CollectionTask:
		return &Task{}, nil
	case CollectionComment:
		return &Comment{}, nil
	default:
		return nil, errors.New("item not supported")
	}
}

// Collection creates a collection corresponding to the collection name
func (s service) Collection(cname string, keys []string, limit, offset int) (Collection, error) {
	col := &collection{
		lim:     limit,
		off:     offset,
		colname: cname,
	}
	switch cname {
	case CollectionUser:
		return nil, nil
	case CollectionGoal, CollectionFollower:
		var u string
		if err := writeKey(keys, &u); err != nil {
			return nil, err
		}
		col.colkeys = map[string]interface{}{
			"user": u,
		}

	case CollectionHabit, CollectionTask, CollectionComment:
		var (
			u = ""
			g bson.ObjectId
		)
		if err := writeKey(keys, &u, &g); err != nil {
			return nil, err
		}
		col.colkeys = map[string]interface{}{
			"user": u,
			"goal": g,
		}
	default:
		return nil, errors.New("collection not supported")
	}
	switch cname {
	case CollectionUser:
		col.items = &[]User{}
	case CollectionGoal:
		col.items = &[]Goal{}
	case CollectionFollower:
		col.items = &[]Follower{}
	case CollectionHabit:
		col.items = &[]Habit{}
	case CollectionTask:
		col.items = &[]Follower{}
	case CollectionComment:
		col.items = &[]Comment{}
	default:
		return nil, errors.New("collection not supported")
	}
	return col, nil
}

// EncryptPassword encrypts the password using the same algorithm applied on the
// previously stored password
func (s service) EncryptPassword(salt string, password string) string {
	return encryptPassword(salt, password)
}

// EncryptPassword encrypts a password using the provided salt
// Same salt and password always result in a same encrypted password (no side
// effects )
func encryptPassword(salt string, password string) string {
	var appended = []byte(password + salt)
	h := sha256.New()
	h.Write(appended)
	return hex.EncodeToString(h.Sum(nil))
}

// New creates a new data service
func New(databaseURL string, databaseName string) (Service, error) {
	// Connect to mongodb
	mongo, err := mgo.DialWithTimeout(databaseURL, 1*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connect to mongodb: %s", err)
	}
	mongo.SetMode(mgo.Monotonic, true)
	if err := mongo.Ping(); err != nil {
		return nil, fmt.Errorf("connect to mongo: ping: %s", err)
	}

	d := &database{
		session:  mongo,
		Database: mongo.DB(databaseName),
	}
	if err := d.setup(); err != nil {
		return nil, fmt.Errorf("set up database: %s", err)
	}
	return service{
		database: d,
	}, nil
}
