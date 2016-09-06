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

package errors

import "fmt"

// Validate represents validation error
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

// NotFound represents the resource not found error
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

// Duplicated represents the duplicated error
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
