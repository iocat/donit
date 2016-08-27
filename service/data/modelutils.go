package data

import (
	"errors"
	"fmt"

	"gopkg.in/mgo.v2/bson"
)

const (
	accessPrivate      = "PRIVATE"
	accessForFollowers = "FOR_FOLLOWERS"
	accessPublic       = "PUBLIC"
)

func validateAccessibility(access string) error {
	switch access {
	case accessPrivate, accessForFollowers, accessPublic:
		return nil
	default:
		return newBadData(fmt.Sprintf("accessibility %s is undefined (only %s, %s, and %s allowed)",
			access, accessPrivate, accessPublic, accessForFollowers))
	}
}

func writeString(id *string) func(string) error {
	return func(k string) error {
		if len(k) == 0 {
			return newBadData(fmt.Sprintf("key is not provided"))
		}
		*id = k
		return nil
	}
}

func writeObjectId(objid *bson.ObjectId) func(string) error {
	return func(data string) error {
		if !bson.IsObjectIdHex(data) {
			return newBadData(fmt.Sprintf("invalid keys: %s is not a 12 byte hex", data))
		}
		*objid = bson.ObjectIdHex(data)
		return nil
	}
}

// writeKey writes the key to the array of pointer
// and optionally write the last field
func writeKey(k []string, keys ...interface{}) error {
	// Check the range of given keys
	if l, c := len(k), len(keys); !(l == c || l == c-1) {
		return fmt.Errorf("inadequate numbers of keys, got %d, expect [%d,%d]", l, c, c+1)
	}
	for i, s := range k {
		switch key := keys[i].(type) {
		case *string:
			if err := writeString(key)(s); err != nil {
				return err
			}
		case *bson.ObjectId:
			if err := writeObjectId(key)(s); err != nil {
				return err
			}
		default:
			panic(errors.New("not valid type"))

		}
	}
	return nil
}
