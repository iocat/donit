package handler

import (
	"net/http"

	"github.com/iocat/donit/handler/internal/utils"
	"github.com/iocat/donit/internal/achieving"
	"github.com/iocat/donit/internal/achieving/validator"
)

// handler will receive a key if getResourceKey is marked true
func decorateAchievableHandler(getResourceKey bool, handler func(achieving.Goal, string, http.ResponseWriter, *http.Request)) http.HandlerFunc {
	var keyGeneratorFunc func() []string
	if getResourceKey {
		keyGeneratorFunc = Achievable.resourceKeyNames
	} else {
		keyGeneratorFunc = Achievable.collectionKeyNames
	}

	// getParentResource gets the parent resource of the Achievable
	var getParentResource = func(username string, id string) (achieving.Goal, error) {
		user, err := store.RetrieveUser(username)
		if err != nil {
			return nil, err
		}
		goal, err := user.RetrieveGoal(id)
		if err != nil {
			return nil, err
		}
		return goal, nil
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ids []string
			err error
		)
		// get the username and the goal id
		ids, err = utils.MuxGetParams(r, keyGeneratorFunc()...)
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		username, goalid := ids[0], ids[1]
		goal, err := getParentResource(username, goalid)
		if err != nil {
			utils.HandleError(err, w)
			return
		}
		// Get the resource id by request
		var achid string
		if getResourceKey {
			achid = ids[2]
		}
		handler(goal, achid, w, r)
	})
}

// CreateAchievable creates an achievable task
var CreateAchievable = decorateAchievableHandler(false, createAchievable)

// UpdateAchievable updates an achievable task
var UpdateAchievable = decorateAchievableHandler(true, updateAchievable)

// DeleteAchievable deletes an achievable task
var DeleteAchievable = decorateAchievableHandler(true, deleteAchievable)

// AllAchievables reads a list of achievable tasks
var AllAchievables = decorateAchievableHandler(false, allAchievables)

func createAchievable(goal achieving.Goal, _ string, w http.ResponseWriter, r *http.Request) {
	ach, err := validator.Validate(r.Body, Achievable.interpreter())
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	_, err = goal.AddAchievable(ach.(achieving.Achievable))
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(nil, w, http.StatusCreated)
}

func updateAchievable(goal achieving.Goal, achid string, w http.ResponseWriter, r *http.Request) {
	ach, err := validator.Validate(r.Body, Achievable.interpreter())
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	err = goal.UpdateAchievable(ach.(achieving.Achievable), achid)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(nil, w, http.StatusOK)
}

func deleteAchievable(goal achieving.Goal, achid string, w http.ResponseWriter, _ *http.Request) {
	err := goal.RemoveAchievable(achid)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(nil, w, http.StatusOK)
}

func allAchievables(goal achieving.Goal, _ string, w http.ResponseWriter, r *http.Request) {
	l, o, err := utils.GetLimitAndOffset(r)
	if err != nil {
		utils.HandleError(err, w)
		return
	}

	achs, err := goal.RetrieveAchievables(l, o)
	if err != nil {
		utils.HandleError(err, w)
		return
	}
	utils.WriteJSONtoHTTP(achs, w, http.StatusOK)
}
