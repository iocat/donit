package data

import "errors"

var (
	errDuplicate = errors.New("the resource is duplicated")
	errNotFound  = errors.New("the resource is not found")
)

// ErrBadData returns whether this is an error caused by bad data
func ErrBadData(err error) bool {
	if _, ok := err.(errBadData); ok {
		return true
	}
	return false
}

// ErrDuplicate occurs when the data is duplicated
func ErrDuplicate(err error) bool {
	return err == errDuplicate
}

// ErrNotFound occurs when data is not found
func ErrNotFound(err error) bool {
	return err == errNotFound
}

type errBadData struct {
	error
}

// newBadData creates a new ErrBadData
func newBadData(err string) error {
	return errBadData{
		error: errors.New(err),
	}
}
