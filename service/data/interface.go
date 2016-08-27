package data

import "gopkg.in/mgo.v2/bson"

const (
	// CollectionUser represents the user collection
	CollectionUser = "users"
	// CollectionGoal represents the goal collection
	CollectionGoal = "goals"
	// CollectionHabit represents the habit collection
	CollectionHabit = "habits"
	// CollectionTask represent the task collection
	CollectionTask = "tasks"
	// CollectionFollower represents the follower collection
	CollectionFollower = "followers"
	// CollectionComment represents a set of comment on the goal
	CollectionComment = "comments"
)

// Validator represents an object that can validate itself before a Write
// Operation can be conducted
type Validator interface {
	// Validate allows the item to have an examination and ensure that data
	// is consistent as expected. The method could mutate the data
	// to fit the business domain
	Validate() error
}

// Item represents an item that can be put into the database
type Item interface {
	// SetKeys sets the item keys using the provided list of id
	// The key order is predefined, look at a specific item for
	// more information on ordering of values
	SetKeys([]string) error

	// keys returns a set of keys: value set representing the item's key
	keys() bson.M
	// cname returns the collection's name of the item
	cname() string
}

// withSubCol allows recursive deletion of subcollections
type withSubCol interface {
	subCols() (bson.M, []string)
}

// iDatabase represents a database that can perform CRUD
type iDatabase interface {
	ReadCollection(col Collection) error
	Read(item Item) error
	Create(item Item) error
	Update(item Item) error
	Delete(item Item) error
}

// Service represents a data service
type Service interface {
	iDatabase

	// Collection allocates a list of item
	Collection(string, []string, int, int) (Collection, error)

	// Item allocates an item using a collection name record
	Item(string) (Item, error)

	// EncryptPassword exposes the algorithm that is used to
	// encrypt the user password
	// EncryptPassword takes a salt and a password (1st and 2nd arg
	// respectively) and produced a hash string
	EncryptPassword(string, string) string
}

// Collection represents a set of items
type Collection interface {
	// Items return a list of item
	Items() interface{}

	// offset gets the offset of the collection
	offset() int
	// limit gets the maximum number of element of collection
	limit() int
	// keys gets the keymap corresponding to this collection
	keys() map[string]interface{}
	// cname gets the collection name of this collection
	cname() string
}
