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

package users

import (
	"encoding/json"
	"testing"

	"github.com/iocat/donit/internal/donitdoc/internal/utils"
)

func TestValidateStoredUser(t *testing.T) {
	cases := []struct {
		data  string
		valid bool
	}{
		{
			data: `{
				"username":"a",
				"email":"felix",
				"defaultAccess":"PRIVATE"
				
			}`,
			valid: false,
		},
		{
			data: `{
				"username":"felixe",
				"email":"feli@hotmail.cat",
				"firstName":"thanh",
				"lastName":"ngo"
				
			}`,
			valid: false,
		}, {
			data: `{
				"username":"fel",
				"email":"felix@gmail.com",
				"firstName":"Thanh",
				"lastName":"Ngo",
				"defaultAccess":"PUBLIC"
			}`,
			valid: true,
		}, {
			data: `{
				"username":"fel",
				"email":"felix@gmail.com",
				"firstName":"Thanh",
				"lastName":"Ngo",
				"defaultAccess":"PUBLIC",
				"password":null
			}`,
			valid: true,
		}, {
			data: `{
				"username":"whateasdasdsadasdasdwqewewqewqeqeqwewqewqewqeqwewqeqwewqeqweqweqwe"
			}`,
			valid: false,
		}, {
			data: `{
				"username":"fel",
				"email":"felix@gmail.com",
				"firstName":"Thanh",
				"lastName":"Ngo",
				"defaultAccess":"PUBLIC",
				"password":"qwlkeqwjelwqjelqwjelwkqje"
			}`,
			valid: false,
		},
	}
	for i, test := range cases {
		var u StoredUser
		err := json.Unmarshal([]byte(test.data), &u)
		if err != nil {
			t.Fatalf("unable to unmarshal the data: %s", err)
		}
		err = utils.Validate(u)
		if (test.valid && err == nil) || (!test.valid && err != nil) {
			t.Logf("SUCCESS: validate case %d, got error: %v", i+1, err)
			continue
		} else if !test.valid {
			t.Fatalf("FAILURE: validate case %d , got nil error, expect an error", i+1)
		} else {
			t.Fatalf("FAILURE: validate case %d, got error %s, expect nil error", i+1, err)
		}

	}
}
