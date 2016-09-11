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

// Package validator contains code that validate JSON object from the achieving
// package.
package validator

import (
	"fmt"
	"io"

	"github.com/iocat/donit/internal/achieving/errors"
	"github.com/iocat/donit/internal/achieving/internal/achievable"
	"github.com/iocat/donit/internal/achieving/internal/goal"
	"github.com/iocat/donit/internal/achieving/internal/user"
	json "github.com/iocat/donit/internal/achieving/jsoninterpreter"
	valid "gopkg.in/asaskevich/govalidator.v4"
)

func init() {
	valid.CustomTypeTagMap.Set("goalAccessField",
		valid.CustomTypeValidator(goal.AccessValidatorFunc))
	valid.CustomTypeTagMap.Set("validateUserStatus",
		valid.CustomTypeValidator(user.ValidateUserStatus))
	valid.CustomTypeTagMap.Set("validateStatus",
		valid.CustomTypeValidator(achievable.ValidateStatus))
	valid.CustomTypeTagMap.Set("cycle",
		valid.CustomTypeValidator(achievable.ValidateCycle))
	valid.CustomTypeTagMap.Set("daysInWeekOrMonth",
		valid.CustomTypeValidator(achievable.ValidateDaysInWeekOrMonth))
}

// Validate decodes the json body and returns an object corresponding to the json
// interpreter
func Validate(r io.Reader, interpreter json.Interpreter) (interface{}, error) {
	// Decode the body
	var (
		obj interface{}
		err error
	)
	if obj, err = interpreter.Decode(r); err != nil {
		return nil, err
	}
	// Run the validator
	err = validate(obj)
	if err != nil {
		return nil, errors.NewValidate(err.Error())
	}
	return obj, nil
}

// Validate dispatches the validator on the object. Returns any error reported from the validator
func validate(obj interface{}) error {
	if ok, err := valid.ValidateStruct(obj); !ok {
		var got error
	loop:
		for {
			switch e := err.(type) {
			case valid.Error:
				got = e
				break loop
			case valid.Errors:
				err = e.Errors()[0]
			default:
				panic(fmt.Errorf("unexpected type %T", e))
			}
		}
		return got
	} else if err != nil {
		return fmt.Errorf("validate %T error: %s", obj, err)
	}
	return nil
}
