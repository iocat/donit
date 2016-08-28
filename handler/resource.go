package handler

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iocat/donit/service/data"
)

// Resource represents a resource which has 2 handlers
type Resource interface {
	item() (data.Item, error)
	collection([]string, int, int) (data.Collection, error)
	getIDs(*http.Request) ([]string, error)

	URL() (string, string)
	Collection() func(w http.ResponseWriter, r *http.Request)
	Item() func(w http.ResponseWriter, r *http.Request)
}

func (res resource) URL() (string, string) {
	url := ""
	for i, src := range res.idSet {
		if i != res.len()-1 {
			url = fmt.Sprintf("%s/%ss/{%s}", url, src, src)
		} else {
			url = fmt.Sprintf("%s/%ss", url, src)
		}
	}
	return url, fmt.Sprintf("%s/{%s}", url, res.idSet[len(res.idSet)-1])
}

type Handlers struct {
	Resources map[string]Resource
	Validator func(http.ResponseWriter, *http.Request)
	dt        data.Service
}

// GetHandlers gets the handlers
func New(databaseURL, dbName string) (*Handlers, error) {
	var err error
	dt, err := data.New(databaseURL, dbName)
	if err != nil {
		return nil, fmt.Errorf("unable to create a data service: %s", err)
	}

	resources := map[string]Resource{
		data.CollectionUser: &resource{
			idSet: []string{"user"},
			cname: data.CollectionUser,
			dt:    dt,
		},
		data.CollectionGoal: &resource{
			idSet: []string{"user", "goal"},
			cname: data.CollectionGoal,
			dt:    dt,
		},
		data.CollectionFollower: &resource{
			idSet: []string{"user", "follower"},
			cname: data.CollectionFollower,
			dt:    dt,
		},
		data.CollectionHabit: &resource{
			idSet: []string{"user", "goal", "habit"},
			cname: data.CollectionHabit,
			dt:    dt,
		},
		data.CollectionTask: &resource{
			idSet: []string{"user", "goal", "task"},
			cname: data.CollectionTask,
			dt:    dt,
		},
		data.CollectionComment: &resource{
			idSet: []string{"user", "goal", "comment"},
			cname: data.CollectionComment,
			dt:    dt,
		},
	}
	h := &Handlers{
		Resources: resources,
		dt:        dt,
	}
	h.Validator = h.validate
	return h, nil
}

type idSet []string

func (is idSet) IDNames() ([]string, string) {
	return is[:is.len()-1], is[is.len()-1]
}

func (is idSet) len() int {
	return len(is)
}

type resource struct {
	dt    data.Service
	cname string
	idSet
}

// item gets an empty item corresponding to this resource
func (res *resource) item() (data.Item, error) {
	item, err := res.dt.Item(res.cname)
	if err != nil {
		return nil, err
	}
	return item, nil
}

// get a readonly data.Collection
func (res *resource) collection(ids []string, limit, offset int) (data.Collection, error) {
	col, err := res.dt.Collection(res.cname, ids, limit, offset)
	if err != nil {
		return nil, err
	}
	return col, nil
}

// getIds gets the parent ids and attempts to get the child id. If the child id
// doesn't exist getIds will skip the child id
func (res *resource) getIDs(r *http.Request) ([]string, error) {
	ids := make([]string, 0, res.len())
	keys := mux.Vars(r)
	for _, k := range res.idSet[0 : res.len()-1] {
		id, ok := keys[k]
		if !ok {
			return nil, newError(codeInternal, fmt.Sprintf("unable to recognize required key: want %s", k))
		}
		ids = append(ids, id)
	}
	id, ok := keys[res.idSet[res.len()-1]]
	if !ok {
		return ids, nil
	}
	ids = append(ids, id)

	return ids, nil
}

func (res *resource) Collection() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ids, err := res.getIDs(r)
		if err != nil {
			handleError(err, w)
			return
		}
		switch r.Method {
		case "GET":
			// Does not allow get method on the user's collection
			if res.cname == data.CollectionUser {
				handleError(errMethodNotAllowed, w)
				return
			}
			// Get limit and offset fields
			off, lim, err := getLimitAndOffset(r)
			if err != nil {
				handleError(err, w)
				return
			}
			// Create the corresponding collection
			col, err := res.collection(ids, lim, off)
			if err != nil {
				handleError(err, w)
				return
			}
			// Read the collection
			err = res.dt.ReadCollection(col)
			if err != nil {
				handleError(err, w)
				return
			}
			writeJSONtoHTTP(col.Items(), w, http.StatusOK)
			return
		case "POST":
			item, err := res.item()
			if err != nil {
				handleError(err, w)
				return
			}
			err = decodeBodyIntoItem(item, r)
			if err != nil {
				handleError(err, w)
				return
			}
			err = item.SetKeys(ids)
			if err != nil {
				handleError(err, w)
				return
			}
			generated, err := res.dt.Create(item)
			if err != nil {
				handleError(err, w)
				return
			}
			var location = r.URL.EscapedPath()
			if generated != nil {
				location = fmt.Sprintf("%s/%s", location, *generated)
				w.Header().Add("Location", location)
			}
			writeJSONtoHTTP(item, w, http.StatusCreated)
		default:
			handleError(errMethodNotAllowed, w)
			return
		}
	}
}

func (res *resource) Item() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		item, err := res.item()
		if err != nil {
			handleError(err, w)
			return
		}
		ids, err := res.getIDs(r)
		if err != nil {
			handleError(err, w)
			return
		}
		err = item.SetKeys(ids)
		if err != nil {
			handleError(err, w)
			return
		}
		switch r.Method {
		case "GET":
			err := res.dt.Read(item)
			if err != nil {
				handleError(err, w)
				return
			}
			// Mask the user's Password and salt
			if res.cname == data.CollectionUser {
				item.(*data.User).Password = nil
				item.(*data.User).Salt = nil
			}
			writeJSONtoHTTP(item, w, http.StatusOK)
			return
		case "PUT":
			err := decodeBodyIntoItem(item, r)
			if err != nil {
				handleError(err, w)
				return
			}
			err = res.dt.Update(item)
			if err != nil {
				handleError(err, w)
				return
			}
		case "DELETE":
			err := res.dt.Delete(item)
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
