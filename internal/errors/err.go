package errors

import "fmt"

type validate struct {
	Field  string
	Reason string
}

// Error implements the error interface
func (v validate) Error() string {
	return fmt.Sprintf("validate field %s: %s", v.Field, v.Reason)
}

// NewValidate returns an Validate error with the field and the reason for
// that error
func NewValidate(field string, reason string) error {
	return &validate{
		Field:  field,
		Reason: reason,
	}
}

// IsValidate returns whether the error is the validation error or not
func IsValidate(err error) bool {
	switch err.(type) {
	case validate, *validate:
		return true
	default:
		return false
	}
}

// notFound represents the resource not found error
type notFound struct {
	ResourceName string
	IdentifiedBy string
}

// Error implements the error interface
func (nf notFound) Error() string {
	return fmt.Sprintf("resource %s identified by %s is not found", nf.ResourceName, nf.IdentifiedBy)
}

// NewNotFound creates a new not found error
func NewNotFound(resource string, key string) error {
	return &notFound{
		ResourceName: resource,
		IdentifiedBy: key,
	}
}

// IsNotFound returns whether the error is not found or not
func IsNotFound(err error) bool {
	switch err.(type) {
	case notFound, *notFound:
		return true
	default:
		return false
	}
}

// duplicated represents the duplicated error
type duplicated struct {
	ResourceName string
	IdentifiedBy string
}

// Error implements the error interface
func (rd duplicated) Error() string {
	return fmt.Sprintf("resource %s identified by %s is duplicated", rd.ResourceName, rd.IdentifiedBy)
}

// NewDuplicated creates a resource duplicated error
func NewDuplicated(name string, key string) error {
	return &duplicated{
		ResourceName: name,
		IdentifiedBy: key,
	}
}

// IsDuplicated returns whether the data has duplicated or not
func IsDuplicated(err error) bool {
	switch err.(type) {
	case duplicated, *duplicated:
		return true
	default:
		return false
	}
}
