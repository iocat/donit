package handler

import (
	"net/http"

	"github.com/iocat/donit/internal/achieving/validator"

	"github.com/iocat/donit/handler/internal/utils"
	"github.com/iocat/donit/internal/achieving"
)

func decorateGoalHandler(getResourceKey bool, handler func(achieving.User, string, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	var keyGeneratorFunc func() []string
	if getResourceKey {
		keyGeneratorFunc = Goal.resourceKeyNames
	} else {
		keyGeneratorFunc = Goal.collectionKeyNames
	}
	var getParentResource = func(username string) (achieving.User, error) {
		user, err := store.RetrieveUser(username)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids, err := utils.MuxGetParams(r, keyGeneratorFunc()...)
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		username := ids[0]
		user, err := getParentResource(username)
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		var gid string
		if getResourceKey {
			gid = ids[1]
		}
		handler(user, gid, w, r)
	})
}

var CreateGoal = decorateGoalHandler(false, createGoal)
var UpdateGoal = decorateGoalHandler(true, updateGoal)
var DeleteGoal = decorateGoalHandler(true, deleteGoal)
var ReadGoal = decorateGoalHandler(true, readGoal)
var AllGoals = decorateGoalHandler(false, allGoals)

func createGoal(user achieving.User, _ string, w http.ResponseWriter, r *http.Request) {
	goal, err := validator.Validate(r.Body, Goal.interpreter())
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	_, err = user.CreateGoal(goal.(achieving.Goal))
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	// TODO: Figure out how to set the resource location
	utils.WriteJSONtoHTTP(nil, w, http.StatusCreated)
}

func updateGoal(user achieving.User, goalid string, w http.ResponseWriter, r *http.Request) {
	goal, err := validator.Validate(r.Body, Goal.interpreter())
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	err = user.UpdateGoal(goal.(achieving.Goal), goalid)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	// TODO: Figure out how to set the resource location
	utils.WriteJSONtoHTTP(nil, w, http.StatusNoContent)
}

func deleteGoal(user achieving.User, gid string, w http.ResponseWriter, r *http.Request) {

	err := user.DeleteGoal(gid)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(nil, w, http.StatusNoContent)
}

func readGoal(user achieving.User, gid string, w http.ResponseWriter, r *http.Request) {
	goal, err := user.RetrieveGoal(gid)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(goal, w, http.StatusOK)
}

func allGoals(user achieving.User, _ string, w http.ResponseWriter, r *http.Request) {
	l, o, err := utils.GetLimitAndOffset(r)
	if err != nil {
		utils.HandleError(err, w)
		return
	}

	gs, err := user.RetrieveGoals(l, o)
	if err != nil {
		utils.HandleError(err, w)
		return
	}

	utils.WriteJSONtoHTTP(gs, w, http.StatusOK)
}
