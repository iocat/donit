package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/iocat/donit/service/data"
)

const (
	codeDecodeJSON = iota + 1
	codeBadForm
	codeMethodNotAllowed
	codeInternal

	codeResourceNotFound
	codeResourceDuplicate
	codeBadData
)

var (
	errNotFound          = newError(codeResourceNotFound, "resource not found")
	errMethodNotAllowed  = newError(codeMethodNotAllowed, "method not allowed")
	errInternal          = newError(codeInternal, "internal server error")
	errBadData           = newError(codeBadData, "the provided data is invalid")
	errBadForm           = newError(codeBadForm, "the provided form is invalid")
	errDecodeJSON        = newError(codeDecodeJSON, "unable to decode the JSON message")
	errDuplicateResource = newError(codeResourceDuplicate, "resource duplicated")
)

// Error represents a handler error
type Error struct {
	Code   int    `json:"code,omitempty"`
	Reason string `json:"reason,omitempty"`
}

func (err Error) Error() string {
	return fmt.Sprintf("ERROR %d: %s", err.Code, err.Reason)
}

func newError(code int, reason interface{}) Error {
	r := ""
	switch reason := reason.(type) {
	case error:
		r = reason.Error()
	case string:
		r = reason
	default:
		panic(errors.New("error type not supported"))
	}
	return Error{
		Code:   code,
		Reason: r,
	}
}

func newInternal(reason interface{}) Error {
	return newError(codeInternal, reason)
}

// NOTE+TODO: normal + non-serialized errors are considered internal error
// TODO: log internal errors
func handleError(err error, w http.ResponseWriter) {
	// convert data service's error
	switch {
	case data.ErrBadData(err):
		err = newError(codeBadData, err)
	case data.ErrDuplicate(err):
		err = newError(codeResourceDuplicate, err)
	case data.ErrNotFound(err):
		err = newError(codeResourceNotFound, err)
	default:
	}

	// handle local package's error
	if err, ok := err.(Error); ok {
		switch err.Code {
		case codeInternal:
			writeJSONtoHTTP(err, w, http.StatusInternalServerError)
			// TODO: log
		case codeResourceNotFound:
			writeJSONtoHTTP(err, w, http.StatusNotFound)
		case codeMethodNotAllowed:
			writeJSONtoHTTP(err, w, http.StatusMethodNotAllowed)
		case codeResourceDuplicate:
			writeJSONtoHTTP(err, w, http.StatusConflict)
		case codeBadData, codeDecodeJSON, codeBadForm:
			writeJSONtoHTTP(err, w, http.StatusBadRequest)
		default:
			errLog.Printf("unrecognized error code: %d", err.Code)
		}
		return
	}
	fmt.Println(err)
	writeJSONtoHTTP(nil, w, http.StatusInternalServerError)

}
