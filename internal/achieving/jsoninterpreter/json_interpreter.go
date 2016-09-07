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

// Package jsoninterpreter handles JSON marshaling and json unmarshaling
package jsoninterpreter

import (
	"fmt"
	"io"

	"gopkg.in/mgo.v2"

	"github.com/iocat/donit/internal/achieving"
)

type errInvalidJSONType string

func (err errInvalidJSONType) Error() string {
	return string(err)
}

// IsErrInvalidJSONType checks whether the error is caused by
// wrong types
func IsErrInvalidJSONType(err error) bool {
	_, ok := err.(errInvalidJSONType)
	return ok
}

func newErrInvalidJSONType(err string) error {
	return errInvalidJSONType(err)
}

// Interpreter represents an interpreter
type Interpreter interface {
	// Decode interprets and decodes the JSON data into the output object
	Decode(io.Reader) (interface{}, error)
	// Encode writes the object to the io.Writer output buffer
	// If the object implements the DataGetter interface{}, the data returned
	// by DataGetter will be written instead, it is encouraged
	// that the object implements the DataGetter interface.
	//
	// If the object implements json.Marshaler, the marshalling output will be
	// written instead
	Encode(io.Writer, interface{}) error
}

// UserJSONInterpreter implements Interpreter
type UserJSONInterpreter struct {
	userCol    *mgo.Collection
	goalCol    *mgo.Collection
	achievable *mgo.Collection
}

func (u UserJSONInterpreter) decode(r io.Reader) (achieving.User, error) {

	return nil, nil
}

func (u UserJSONInterpreter) Decode(r io.Reader) (interface{}, error) {

	return u.decode(r)
}

func (u UserJSONInterpreter) encode(w io.Writer, user achieving.User) error {

	return nil
}

func (u UserJSONInterpreter) Encode(w io.Writer, user interface{}) error {
	casted, ok := user.(achieving.User)
	if !ok {
		return newErrInvalidJSONType(fmt.Sprintf("the provided type is not User, got %T", user))
	}
	return u.encode(w, casted)
}
