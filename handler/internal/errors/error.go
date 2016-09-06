package errors

import (
	stderr "errors"
	"fmt"
	"net/http"

	docerr "github.com/iocat/donit/internal/donitdoc/errors"
)

const (
	codeDecodeJSON code = iota + 1
	codeMethodNotAllowed
	codeInternal
	codeResourceNotFound
	codeResourceDuplicate
	codeBadData
)

type code int

// HTTPStatus returns the http status code associated with the error
func (c code) HTTPStatus() int {
	switch c {
	case codeInternal:
		return http.StatusInternalServerError
	case codeResourceNotFound:
		return http.StatusNotFound
	case codeMethodNotAllowed:
		return http.StatusMethodNotAllowed
	case codeResourceDuplicate:
		return http.StatusConflict
	case codeBadData, codeDecodeJSON:
		return http.StatusBadRequest
	default:
		panic(stderr.New("unknown http status code for this error"))
	}
}

var (
	// ErrNotFound ...
	ErrNotFound = newError(codeResourceNotFound, "resource not found")
	// ErrMethodNotAllowed ...
	ErrMethodNotAllowed = newError(codeMethodNotAllowed, "method not allowed")
	// ErrInternal ...
	ErrInternal = newError(codeInternal, "internal server error")
	// ErrBadData ...
	ErrBadData = newError(codeBadData, "the provided data is invalid")
	// ErrDecodeJSON ...
	ErrDecodeJSON = newError(codeDecodeJSON, "unable to decode the JSON message")
	// ErrDuplicateResource ...
	ErrDuplicateResource = newError(codeResourceDuplicate, "resource duplicated")
)

// Error represents a handler error
type Error struct {
	Status string `json:"_status"`
	Code   code   `json:"code,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func (err Error) Error() string {
	return fmt.Sprintf("ERROR %d: %s", err.Code, err.Reason)
}

// NewBadData creates a new bad data error object
func NewBadData(reason interface{}) error {
	return newError(codeBadData, reason)
}

// NewInternal creates a new internal error object
func NewInternal(reason interface{}) error {
	return newError(codeInternal, reason)
}

func newError(c code, reason interface{}) Error {
	r := ""
	switch reason := reason.(type) {
	case error:
		r = reason.Error()
	case string:
		r = reason
	default:
		panic(stderr.New("error type not supported"))
	}
	return Error{
		Status: "error",
		Code:   c,
		Reason: r,
	}
}

// ParseDocumentError converts the known document error into a serialized error
func ParseDocumentError(err error) error {
	switch {
	case docerr.IsDuplicated(err):
		return newError(codeResourceDuplicate, err)
	case docerr.IsValidate(err):
		return newError(codeBadData, err)
	case docerr.IsNotFound(err):
		return newError(codeResourceNotFound, err)
	default:
		return err

	}
}
