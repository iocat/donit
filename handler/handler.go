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
	parseIDToItem(&user, r)
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

func decodeBodyIntoItem(obj Item, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		return errDecodeJSON
	}
	return nil
}

func getCollectionIDs(collection string, r *http.Request) bson.M {
	m := make(bson.M)
	keys := mux.Vars(r)
	switch collection {
	case collectionUser:
		return nil
	case collectionGoal:
		m["user"] = keys["user"]
	case collectionHabit, collectionTask:
		m["user"] = keys["user"]
		m["goal"] = keys["goal"]
	default:
		panic(errors.New("type not supported"))
	}
	return nil
}

func getCollection(collection string) interface{} {
	switch collection {
	case collectionUser:
		return nil
	case collectionGoal:
		return []Goal{}
	case collectionHabit:
		return []Habit{}
	case collectionTask:
		return []Task{}
	default:
		panic("collection not supported")
	}
}

func parseIDToItem(obj Item, r *http.Request) {
	keys := mux.Vars(r)
	switch obj := obj.(type) {
	case *User:
		obj.Username = keys["user"]
	case *Goal:
		obj.User = keys["user"]
		obj.ID = keys["goal"]
	case *Habit:
		obj.User = keys["user"]
		obj.Goal = keys["goal"]
		obj.ID = keys["habit"]
	case *Task:
		obj.User = keys["user"]
		obj.Goal = keys["goal"]
		obj.ID = keys["task"]
	default:
		panic("type not supported")
	}
}

func getItem(collection string) Item {
	switch collection {
	case collectionUser:
		return &User{}
	case collectionGoal:
		return &Goal{}
	case collectionHabit:
		return &Habit{}
	case collectionTask:
		return &Task{}
	default:
		panic("collection not supported")
	}
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
			id, col := getCollectionIDs(c, r), getCollection(c)
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
		case "POST":
			err := decodeBodyIntoItem(item, r)
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
		parseIDToItem(item, r)
		switch r.Method {
		case "GET":
			err := db.Read(item)
			if err != nil {
				handleError(err, w)
				return
			}
			// Mask the user's Password
			if c == collectionUser {
				item.(*User).Password = nil
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
