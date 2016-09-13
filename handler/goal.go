package handler

import (
	"net/http"

	"github.com/iocat/donit/internal/achieving/validator"

	"github.com/iocat/donit/handler/internal/utils"
	"github.com/iocat/donit/internal/achieving"
)

func CreateGoal(w http.ResponseWriter, r *http.Request) {
	ids, err := utils.MuxGetParams(r, Goal.collectionKeyNames()...)
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
	goal, err := validator.Validate(r.Body, Goal.interpreter())
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	id, err := user.CreateGoal(goal.(achieving.Goal))
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	// TODO: Figure out how to set the resource location
	utils.WriteJSONtoHTTP(nil, w, http.StatusCreated)
}

func UpdateGoal(w http.ResponseWriter, r *http.Request) {
	ids, err := utils.MuxGetParams(r, Goal.resourceKeyNames()...)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	username, goalid := ids[0], ids[1]
	user, err := store.RetrieveUser(username)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	goal, err := validator.Validate(r.Body, Goal.interpreter())
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	gid, err := achieving.CreateID(goalid)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	err = user.UpdateGoal(goal.(achieving.Goal), gid)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	// TODO: Figure out how to set the resource location
	utils.WriteJSONtoHTTP(nil, w, http.StatusNoContent)
}

func DeleteGoal(w http.ResponseWriter, r *http.Request) {
	ids, err := utils.MuxGetParams(r, Goal.resourceKeyNames()...)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	username, goalid := ids[0], ids[1]
	gid, err := achieving.CreateID(goalid)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	user, err := store.RetrieveUser(username)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	err = user.DeleteGoal(gid)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(nil, w, http.StatusNoContent)
}

func ReadGoal(w http.ResponseWriter, r *http.Request) {
	ids, err := utils.MuxGetParams(r, Goal.resourceKeyNames()...)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	username, goalid := ids[0], ids[1]
	gid, err := achieving.CreateID(goalid)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	user, err := store.RetrieveUser(username)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	goal, err := user.RetrieveGoal(gid)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(goal, w, http.StatusOK)
}
