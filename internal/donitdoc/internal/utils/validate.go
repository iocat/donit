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

package utils

import (
	"fmt"

	valid "gopkg.in/asaskevich/govalidator.v4"
)

// Validate dispatches the validator on the object. Returns any error reported from the validator
func Validate(obj interface{}) error {
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
