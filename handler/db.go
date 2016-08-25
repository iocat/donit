package handler

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// newDatabase creates a new database
func newDatabase() (*database, error) {
	var defaultDatabaseURL = "127.0.0.1"
	var defaultDatabaseName = "donit"

	mongo, err := mgo.DialWithTimeout(defaultDatabaseURL, 1*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connect to mongodb: %s", err)
	}
	mongo.SetMode(mgo.Monotonic, true)
	if err := mongo.Ping(); err != nil {
		return nil, fmt.Errorf("connect to mongo: ping: %s", err)
	}

	d := &database{
		session:  mongo,
		Database: mongo.DB(defaultDatabaseName),
	}
	if err := d.setup(); err != nil {
		return nil, fmt.Errorf("set up database: %s", err)
	}
	return d, nil
}

type database struct {
	session *mgo.Session
	*mgo.Database
}

func (d *database) setup() error {
	// Build the index
	usernameIndex := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}
	// Make sure user database has a secondary index
	if err := d.C(collectionUser).EnsureIndex(usernameIndex); err != nil {
		return err
	}
	return nil
}

// Item represents an item that can be put into the database
type Item interface {
	// KeySet returns a key : value set representing the item's key
	KeySet() bson.M

	// Cname returns the collection's name of the item
	Cname() string
}

// Validator represents an object that can validate itself before a Write
// Operation can be conducted
// Validate allows the item to make an examination and ensure that the data
// is consistent as expected. It may change the data to fit the business
// domain requirement
type Validator interface {
	Validate() error
}

// If limit == 0, limit is not apploed, offset is always applied
// keys are the set of keys used to retrieve the item
// res is the result set
func (db database) GetCollection(collection string, offset, limit int, res interface{}, keys bson.M) error {
	q := db.C(collection).Find(keys)
	if limit > 0 {
		q = q.Limit(limit)
	}
	if offset >= 0 {
		q = q.Skip(offset)
	}
	err := q.All(&res)
	if err != nil {
		return err
	}
	return nil
}

// ReadItem reads data into the item
func (db database) Read(item Item) error {
	if err := db.C(item.Cname()).Find(item.KeySet()).One(item); err != nil {
		if err == mgo.ErrNotFound {
			return errNotFound
		}
		return err
	}
	return nil
}

// CreateItem creates a new item
func (db database) Create(item Item) error {
	if v, ok := item.(Validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	if err := db.C(item.Cname()).Insert(item); err != nil {
		if mgo.IsDup(err) {
			return errDuplicateResource
		}
		return err
	}
	return nil
}

// DeleteItem deletes an item
func (db database) Delete(item Item) error {
	if _, err := db.C(item.Cname()).RemoveAll(item.KeySet()); err != nil {
		if err == mgo.ErrNotFound {
			return errNotFound
		}
		return err
	}
	return nil
}

// UpdateItem updates an item
func (db database) Update(item Item) error {
	if v, ok := item.(Validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	if err := db.C(item.Cname()).Update(item.KeySet(), item); err != nil {
		if err == mgo.ErrNotFound {
			return errNotFound
		}
		return err
	}
	return nil
}
