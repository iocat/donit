package data

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type database struct {
	session *mgo.Session
	*mgo.Database
}

func (db *database) setup() error {
	// Build the index
	usernameIndex := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}
	// Make sure user database has a secondary index
	if err := db.C(CollectionUser).EnsureIndex(usernameIndex); err != nil {
		return err
	}
	// follower index
	followerIndex := mgo.Index{
		Key:        []string{"username", "follower"},
		Unique:     true,
		DropDups:   true,
		Background: true,
	}
	if err := db.C(CollectionFollower).EnsureIndex(followerIndex); err != nil {
		return err
	}
	return nil
}

//  GetCollection gets the collection and put it in res
// If limit == 0, limit is not apploed, offset is always applied
// keys are the set of keys used to retrieve the item
// res is the result set
func (db database) ReadCollection(c Collection) error {
	q := db.C(c.cname()).Find(c.keys())
	if c.limit() > 0 {
		q = q.Limit(c.limit())
	}
	if c.offset() >= 0 {
		q = q.Skip(c.offset())
	}
	err := q.All(c.Items())
	if err != nil {
		return err
	}
	return nil
}

// Read reads data into the item
func (db database) Read(item Item) error {
	if err := db.C(item.cname()).Find(item.keys()).One(item); err != nil {
		if err == mgo.ErrNotFound {
			return errNotFound
		}
		return err
	}
	return nil
}

// Create creates a new item. Create also updates the object value in case new
// id is generated. The first returned argument signals a new id is created
func (db database) Create(item Item) (*string, error) {
	if v, ok := item.(Validator); ok {
		if err := v.Validate(); err != nil {
			return nil, err
		}
	}
	var generated *string
	if item, ok := item.(IDGenerator); ok {
		generated = new(string)
		*generated = item.GenerateID()
	}
	if err := db.C(item.cname()).Insert(item); err != nil {
		if mgo.IsDup(err) {
			return generated, errDuplicate
		}
		return generated, err
	}
	return generated, nil
}

func (db database) delete(col string, keys bson.M) error {
	if _, err := db.C(col).RemoveAll(keys); err != nil {
		if err == mgo.ErrNotFound {
			return errNotFound
		}
		return err
	}
	return nil
}

// Delete deletes an item
func (db database) Delete(item Item) error {
	ic, k := item.cname(), item.keys()
	if err := db.delete(ic, k); err != nil {
		return err
	}
	if item, ok := item.(withSubCol); ok {
		k, cols := item.subCols()
		for _, c := range cols {
			err := db.delete(c, k)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Update updates an item
func (db database) Update(item Item) error {
	if v, ok := item.(Validator); ok {
		if err := v.Validate(); err != nil {
			return err
		}
	}
	if err := db.C(item.cname()).Update(item.keys(), item); err != nil {
		if err == mgo.ErrNotFound {
			return errNotFound
		}
		return err
	}
	return nil
}
