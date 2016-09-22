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

package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iocat/donit/errors"
)

// WriteJSONtoHTTP writes the object to the HTTP response with the provided http code
// If no object is provided, the content of the response would be empty
// TODO: log JSON error
func WriteJSONtoHTTP(obj interface{}, w http.ResponseWriter, c int) {
	const jsonContentType = "application/json; charset=utf-8"
	// writeJSON writes the data to the output buffer
	var writeJSON = func(obj interface{}, w io.Writer) error {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "\t")
		if err := enc.Encode(obj); err != nil {
			return fmt.Errorf("write JSON to HTTP response: %s", err)
		}
		return nil
	}
	w.Header().Set("Content-Type", jsonContentType)
	w.WriteHeader(c)
	if obj == nil {
		return
	}
	// NOTE: ignore error here
	// (*￣(ｴ)￣*) <- here's a cute bear, just in case you get angry
	_ = writeJSON(obj, w)
}

// TODO: implement this
func WriteJSONtoHTTPWithLocation(loc string, obj interface{}, w http.ResponseWriter, c int) {
	w.Header().Set("Location", loc)
	WriteJSONtoHTTP(obj, w, c)
}

// HandleError is an utility function to handle the error
func HandleError(err error, w http.ResponseWriter) {
	err = errors.ParseDocumentError(err)
	// handle local package's error
	if err, ok := err.(errors.Error); ok {
		// NOTE: Temporarily write to stdout
		fmt.Println(err)
		WriteJSONtoHTTP(err, w, err.Code.HTTPStatus())
		return
	}
	WriteJSONtoHTTP(nil, w, http.StatusInternalServerError)
}

// DecodeJSON reads the Reader and reflects the value into
// the read object
func DecodeJSON(r io.Reader, obj interface{}) error {
	if err := json.NewDecoder(r).Decode(obj); err != nil {
		return errors.ErrDecodeJSON
	}
	return nil
}

// GetLimitAndOffset gets the limit and the offset form values
func GetLimitAndOffset(r *http.Request) (int, int, error) {
	var offs, lim int
	var err error
	if err = r.ParseForm(); err != nil {
		return -1, -1, errors.ErrInternal
	}
	stro := r.Form.Get("offset")
	if len(stro) == 0 {
		offs = 0
	} else if offs, err = strconv.Atoi(stro); err != nil {
		return -1, -1, errors.ErrBadData
	}
	strl := r.Form.Get("limit")
	if len(strl) == 0 {
		lim = -1
	} else if lim, err = strconv.Atoi(strl); err != nil {
		return -1, -1, errors.ErrBadData
	}
	return offs, lim, nil
}

// MuxGetParams gets the request's parameter from the HTTP request's URL
func MuxGetParams(r *http.Request, params ...string) ([]string, error) {
	v := mux.Vars(r)
	var res = make([]string, 0, len(params))
	for _, p := range params {
		r, ok := v[p]
		if !ok {
			return nil, errors.NewInternal("cannot find id for " + p)
		}
		res = append(res, r)
	}
	return res, nil
}
