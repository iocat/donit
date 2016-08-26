package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

var db *database

func init() {
	var err error
	db, err = newDatabase()
	if err != nil {
		errLog.Println(err)
		os.Exit(1)
	}
}

// Get gets the handler corresponding to the handler name
func Get(name string) http.HandlerFunc {
	users, user := generateHandlers(collectionUser)
	goals, goal := generateHandlers(collectionGoal)
	habits, habit := generateHandlers(collectionHabit)
	tasks, task := generateHandlers(collectionTask)
	followers, follower := generateHandlers(collectionFollower)
	switch name {
	case "users":
		return users
	case "user":
		return user
	case "goals":
		return goals
	case "goal":
		return goal
	case "habits":
		return habits
	case "habit":
		return habit
	case "tasks":
		return tasks
	case "task":
		return task
	case "followers":
		return followers
	case "follower":
		return follower
	case "validator":
		return validate
	default:
		panic(fmt.Errorf("unsupported controller: %s", name))
	}
}

func validate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handleError(errBadForm, w)
		return
	}
	password := r.Form.Get("password")
	if len(password) == 0 {
		handleError(newError(codeBadForm, "the password field is not provided"), w)
		return
	}
	user := User{}
	parseIDToItem(&user, true, r)
	err := db.Read(&user)
	if err != nil {
		handleError(err, w)
		return
	}
	if *user.Password == encryptPassword(*user.Salt, password) {
		writeJSONtoHTTP(true, w, http.StatusOK)
		return
	}
	writeJSONtoHTTP(false, w, http.StatusOK)
	return
}

// decodeBodyIntoItem reads the request body and reflects the value into
// the object
func decodeBodyIntoItem(obj Item, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		return errDecodeJSON
	}
	return nil
}

// getCollectionParentIDs get the ID of the parent
func getCollectionParentIDs(collection string, r *http.Request) bson.M {
	m := make(bson.M)
	keys := mux.Vars(r)
	switch collection {
	case collectionUser:
		return nil
	case collectionGoal, collectionFollower:
		m["user"] = keys["user"]
	case collectionHabit, collectionTask:
		m["user"], m["goal"] = keys["user"], keys["goal"]
	default:
		panic(errors.New("type not supported"))
	}
	return nil
}

// getCollection gets the collection concrete type based on the collection's
// name
func getCollection(collection string) interface{} {
	switch collection {
	case collectionUser:
		return nil
	case collectionGoal:
		return &[]Goal{}
	case collectionHabit:
		return &[]Habit{}
	case collectionTask:
		return &[]Task{}
	case collectionFollower:
		return &[]Follower{}
	default:
		panic("collection not supported")
	}
}

// getItem gets the item based on the collection name
func getItem(collection string) Item {
	switch collection {
	case collectionUser:
		return &User{}
	case collectionGoal:
		return &Goal{}
	case collectionFollower:
		return &Follower{}
	case collectionHabit:
		return &Habit{}
	case collectionTask:
		return &Task{}
	default:
		panic("collection not supported")
	}
}

// parseIDToItem reads the parent+child id from the request and reflect those
// id into the item struct
// readChildId indicates whether to read the child's id or not,
// (i.e. false in case we create an auto genereated id, the id provided is unnecessary)
// Error occurs when the item's id is invalid ( not a 12 byte hex string)
func parseIDToItem(obj Item, readChildID bool, r *http.Request) error {
	keys := mux.Vars(r)
	var (
		idName string
		objID  *bson.ObjectId
	)

	switch obj := obj.(type) {
	case *User:
		if readChildID {
			obj.Username = keys["user"]
		}
		return nil
	case *Goal:
		obj.User = keys["user"]
		idName, objID = "goal", &obj.ID
	case *Follower:
		obj.User = keys["user"]
		if readChildID {
			obj.Follower = keys["follower"]
		}
		return nil
	case *Habit:
		obj.User = keys["user"]
		obj.Goal = keys["goal"]
		idName, objID = "habit", &obj.ID
	case *Task:
		obj.User = keys["user"]
		obj.Goal = keys["goal"]
		idName, objID = "task", &obj.ID
	default:
		panic("type not supported")
	}
	if readChildID {
		id := keys[idName]
		if !bson.IsObjectIdHex(id) {
			return newError(codeBadData, fmt.Sprintf("id %s is invalid (must be a 12 byte hex string)", id))
		}
		*objID = bson.ObjectIdHex(id)
	}
	return nil
}

// getLimitAndOffset gets the limit and the offset form values
func getLimitAndOffset(r *http.Request) (int, int, error) {
	var offs, lim int
	var err error
	if err = r.ParseForm(); err != nil {
		return 0, 0, errInternal
	}
	stro := r.Form.Get("offset")
	if len(stro) == 0 {
		offs = 0
	} else if offs, err = strconv.Atoi(stro); err != nil {
		return 0, 0, errBadForm
	}
	strl := r.Form.Get("limit")
	if len(strl) == 0 {
		lim = -1
	} else if lim, err = strconv.Atoi(strl); err != nil {
		return 0, 0, errBadForm
	}
	return offs, lim, nil
}

// createRestHandler creates 2 handler one for the collection (the set without id)
// and one for an item
func generateHandlers(c string) (func(w http.ResponseWriter, r *http.Request), func(w http.ResponseWriter, r *http.Request)) {
	return generateCollectionHandler(c), generateItemHandler(c)
}

func generateCollectionHandler(c string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		item := getItem(c)
		switch r.Method {
		case "GET":
			if c == collectionUser {
				handleError(errMethodNotAllowed, w)
				return
			}
			id, col := getCollectionParentIDs(c, r), getCollection(c)
			off, lim, err := getLimitAndOffset(r)
			if err != nil {
				handleError(err, w)
				return
			}
			err = db.GetCollection(c, off, lim, col, id)
			if err != nil {
				handleError(err, w)
				return
			}
			writeJSONtoHTTP(col, w, http.StatusOK)
			return
		case "POST":
			err := decodeBodyIntoItem(item, r)
			if err != nil {
				handleError(err, w)
				return
			}
			err = parseIDToItem(item, false, r)
			if err != nil {
				handleError(err, w)
				return
			}
			err = db.Create(item)
			if err != nil {
				handleError(err, w)
				return
			}
		default:
			handleError(errMethodNotAllowed, w)
			return
		}
	}
}

func generateItemHandler(c string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		item := getItem(c)
		if err := parseIDToItem(item, true, r); err != nil {
			handleError(err, w)
			return
		}
		switch r.Method {
		case "GET":
			err := db.Read(item)
			if err != nil {
				handleError(err, w)
				return
			}
			// Mask the user's Password and salt
			if c == collectionUser {
				item.(*User).Password = nil
				item.(*User).Salt = nil
			}
			writeJSONtoHTTP(item, w, http.StatusOK)
			return
		case "PUT":
			err := decodeBodyIntoItem(item, r)
			if err != nil {
				handleError(err, w)
				return
			}
			err = db.Update(item)
			if err != nil {
				handleError(err, w)
				return
			}
		case "DELETE":
			err := db.Delete(item)
			if err != nil {
				handleError(err, w)
				return
			}
		default:
			handleError(errMethodNotAllowed, w)
			return
		}
	}
}
