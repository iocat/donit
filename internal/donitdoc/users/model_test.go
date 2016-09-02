package users

import (
	"encoding/json"
	"testing"

	"github.com/iocat/donit/internal/donitdoc/utils"
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
		var u storedUser
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
