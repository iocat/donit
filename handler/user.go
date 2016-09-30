package handler

import (
	"fmt"
	"net/http"
	"path"

	"github.com/iocat/donit/errors"
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

func decorateUserHandler(getResourceKey bool, handler func(achieving.UserStore, string, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	var keyGeneratorFunc func() []string
	if getResourceKey {
		keyGeneratorFunc = User.resourceKeyNames
	} else {
		keyGeneratorFunc = User.collectionKeyNames
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids, err := utils.MuxGetParams(r, keyGeneratorFunc()...)
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		var username string
		if getResourceKey {
			username = ids[0]
		}
		handler(store, username, w, r)
	})
}

var CreateUser = decorateUserHandler(false, createUser)
var ReadUser = decorateUserHandler(true, readUser)
var DeleteUser = decorateUserHandler(true, deleteUser)
var UpdateUser = decorateUserHandler(true, updateUser)
var Auth = decorateUserHandler(true, authUser)
var PasswordChange = decorateUserHandler(true, changePassword)

func changePassword(store achieving.UserStore, username string, w http.ResponseWriter, r *http.Request) {
	getChangePassword := func(param string, r *http.Request) (string, error) {
		if err := r.ParseForm(); err != nil {
			return "", errors.ErrBadData
		}
		password := r.Form.Get(param)
		if len(password) == 0 {
			return "", errors.NewBadData(fmt.Sprintf("password param %s not provided", param))
		}
		return password, nil
	}
	// get two passwords
	oldpass, err := getChangePassword("old", r)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	newpass, err := getChangePassword("new", r)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	// Change the password
	err = store.ChangePassword(username, oldpass, newpass)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(nil, w, http.StatusNoContent)
}

func authUser(store achieving.UserStore, username string, w http.ResponseWriter, r *http.Request) {
	password, err := getPassword(r)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	ok, err := store.Authenticate(username, password)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(ok, w, http.StatusOK)
}

// createUser creates a new user
func createUser(store achieving.UserStore, _ string, w http.ResponseWriter, r *http.Request) {
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
	username, err := store.CreateNewUser(usr, password)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTPWithLocation(path.Join(r.URL.EscapedPath(), username), nil, w, http.StatusCreated)
}

// readUser reads the user data
func readUser(store achieving.UserStore, username string, w http.ResponseWriter, r *http.Request) {
	user, err := store.RetrieveUser(username)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	_ = User.interpreter().Encode(w, user)
}

// DeleteUser deletes an user
func deleteUser(store achieving.UserStore, username string, w http.ResponseWriter, r *http.Request) {
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
func updateUser(store achieving.UserStore, username string, w http.ResponseWriter, r *http.Request) {
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
