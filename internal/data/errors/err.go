package errors

import "fmt"

type Validate struct {
	Reason string
}

// Error implements the error interface
func (v Validate) Error() string {
	return fmt.Sprintf("validate %s", v.Reason)
}

// NewValidate returns an Validate error with the field and the reason for
// that error
func NewValidate(reason string) error {
	return &Validate{
		Reason: reason,
	}
}

// IsValidate returns whether the error is the validation error or not
func IsValidate(err error) bool {
	switch err.(type) {
	case Validate, *Validate:
		return true
	default:
		return false
	}
}

// notFound represents the resource not found error
type NotFound struct {
	ResourceName string
	IdentifiedBy string
}

// Error implements the error interface
func (nf NotFound) Error() string {
	return fmt.Sprintf("resource %s identified by %s is not found", nf.ResourceName, nf.IdentifiedBy)
}

// NewNotFound creates a new not found error
func NewNotFound(resource string, key string) error {
	return &NotFound{
		ResourceName: resource,
		IdentifiedBy: key,
	}
}

// IsNotFound returns whether the error is not found or not
func IsNotFound(err error) bool {
	switch err.(type) {
	case NotFound, *NotFound:
		return true
	default:
		return false
	}
}

// duplicated represents the duplicated error
type Duplicated struct {
	ResourceName string
	IdentifiedBy string
}

// Error implements the error interface
func (rd Duplicated) Error() string {
	return fmt.Sprintf("resource %s identified by %s is duplicated", rd.ResourceName, rd.IdentifiedBy)
}

// NewDuplicated creates a resource duplicated error
func NewDuplicated(name string, key string) error {
	return &Duplicated{
		ResourceName: name,
		IdentifiedBy: key,
	}
}

// IsDuplicated returns whether the data has duplicated or not
func IsDuplicated(err error) bool {
	switch err.(type) {
	case Duplicated, *Duplicated:
		return true
	default:
		return false
	}
}
