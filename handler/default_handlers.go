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

package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/iocat/donit/errors"
	"github.com/iocat/donit/handler/internal/utils"
)

var errLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)

var (
	// NotFound is a default handler for non-supported method or path
	NotFound = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			utils.WriteJSONtoHTTP(errors.ErrNotFound, w, http.StatusNotFound)
		},
	)

	// NotAllowed is a default handler for not-allowed method
	NotAllowed = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			utils.WriteJSONtoHTTP(errors.ErrMethodNotAllowed, w, http.StatusNotFound)
		},
	)
)
