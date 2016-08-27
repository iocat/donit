package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/iocat/donit/service/data"
)

// Get gets the handler which does not deal directly with CRUD resources
func Get(name string) http.HandlerFunc {
	switch name {
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

	ids, err := resources[data.CollectionUser].getIDs(r)
	if err != nil {
		handleError(err, w)
		return
	}
	user := data.User{}
	if err = user.SetKeys(ids); err != nil {
		if err != nil {
			handleError(err, w)
			return
		}
	}
	err = dt.Read(&user)
	if err != nil {
		handleError(err, w)
		return
	}
	if *user.Password == dt.EncryptPassword(*user.Salt, password) {
		writeJSONtoHTTP(true, w, http.StatusOK)
		return
	}
	writeJSONtoHTTP(false, w, http.StatusOK)
	return
}

// decodeBodyIntoItem reads the request body and reflects the value into
// the object
func decodeBodyIntoItem(obj data.Item, r *http.Request) error {
	if err := json.NewDecoder(r.Body).Decode(obj); err != nil {
		return errDecodeJSON
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
