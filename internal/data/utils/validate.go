package utils

import (
	"fmt"

	"github.com/iocat/donit/internal/data/errors"
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
		return errors.NewValidate(got.Error())
	} else if err != nil {
		return fmt.Errorf("validate %T error: %s", obj, err)
	}
	return nil
}
