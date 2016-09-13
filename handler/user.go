package handler

import (
	"net/http"

	"github.com/iocat/donit/handler/internal/errors"
	"github.com/iocat/donit/handler/internal/utils"
	"github.com/iocat/donit/internal/achieving"
	"github.com/iocat/donit/internal/achieving/validator"
)

func getPassword(r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		return "", errors.ErrBadData
	}
	password := r.Form.Get("password")
	if len(password) == 0 {
		return "", errors.NewBadData("password not provided")
	}
	return password, nil
}

// CreateUser creates a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	obj, err := validator.Validate(r.Body, User.interpreter())
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	usr := obj.(achieving.User)
	password, err := getPassword(r)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	err = store.CreateNewUser(usr, password)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(nil, w, http.StatusCreated)
}

// ReadUser reads the user data
func ReadUser(w http.ResponseWriter, r *http.Request) {
	ids, err := utils.MuxGetParams(r, User.resourceKeyNames()...)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	username := ids[0]
	user, err := store.RetrieveUser(username)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	_ = User.interpreter().Encode(w, user)
}

// DeleteUser deletes an user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	ids, err := utils.MuxGetParams(r, User.resourceKeyNames()...)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	username := ids[0]
	password, err := getPassword(r)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	err = store.DeleteUser(username, password)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(nil, w, http.StatusNoContent)
}

// UpdateUser updates an user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	ids, err := utils.MuxGetParams(r, User.resourceKeyNames()...)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	username := ids[0]
	obj, err := validator.Validate(r.Body, User.interpreter())
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	usr := obj.(achieving.User)
	err = store.UpdateUser(usr, username)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(nil, w, http.StatusNoContent)
}
