// Copyright 2016 Thanh Ngo <felix.infinite@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package validator contains code that validate object from the achieving
// package. Validator only validates POST+PUT request with a JSON body.
// User of validator MUST set up the validator properly to
// validate a correct struct.
//
// Validator will NOT propagate the result if the validation process returns an error
// Caller can register a callback function to show a proper response (RegisterInvalidRequestResponse)
package validator

import (
	"net/http"

	"github.com/iocat/donit/internal/achieving/internal/achievable"
	"github.com/iocat/donit/internal/achieving/internal/goal"
	"github.com/iocat/donit/internal/achieving/internal/user"
	valid "gopkg.in/asaskevich/govalidator.v4"
)

func init() {
	valid.SetFieldsRequiredByDefault(true)
	valid.CustomTypeTagMap.Set("goalAccessField", valid.CustomTypeValidator(goal.AccessValidatorFunc))
	valid.CustomTypeTagMap.Set("validateUserStatus", valid.CustomTypeValidator(user.ValidateUserStatus))
	valid.CustomTypeTagMap.Set("validateStatus", valid.CustomTypeValidator(achievable.ValidateStatus))
	valid.CustomTypeTagMap.Set("cycle", valid.CustomTypeValidator(achievable.ValidateCycle))
	valid.CustomTypeTagMap.Set("daysInWeekOrMonth", valid.CustomTypeValidator(achievable.ValidateDaysInWeekOrMonth))
}

// User validates the request body corresponding to an entity, if it is a valid entity
// the method pass it to the http context as a "resource" object.
// ( Caller can use GetValidatedResource helper function to get the "resource"
// value from the request )
func User(handler func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

		})
}

var invalidRequest func(error, http.ResponseWriter)

// RegisterInvalidRequestResponse registers a callback function to handle invalid result
func RegisterInvalidRequestResponse(fn func(error, http.ResponseWriter)) {
	invalidRequest = fn
}

// GetValidatedResource gets the resource corresponding to the request.
func GetValidatedResource(r *http.Request) interface{} {
	return r.Context().Value("resource")
}
