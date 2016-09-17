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

func decorateUserHandler(getResourceKey bool, handler func(achieving.UserStore, string, w http.ResponseWriter, r *http.Request)) http.HandlerFunc{
	var keyGeneratorFunc func() []string
	if getResourceKey {
		keyGeneratorFunc = User.resourceKeyNames
	} else {
		keyGeneratorFunc = User.collectionKeyNames
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		ids,err := utils.MuxGetParams(r, keyGeneratorFunc()...)
		if err != nil{
			utils.HandleError(err,w)
			return
		}
		var username string
		if getResourceKey {
			username = ids[0]
		}
		handler(store, username, w,r)
	})
}

var CreateUser = decorateUserHandler(false, createUser)
var ReadUser = decorateUserHandler(true, readUser)
var DeleteUser = decorateUserHandler(true, deleteUser)
var UpdateUser = decorateUserHandler(true, updateUser)

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
	err = store.CreateNewUser(usr, password)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(nil, w, http.StatusCreated)
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
