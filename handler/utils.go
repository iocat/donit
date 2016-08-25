package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var errLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)

var (
	// NotFound is a default handler for non-supported method or path
	NotFound = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			writeJSONtoHTTP(errNotFound, w, http.StatusNotFound)
		},
	)

	// NotAllowed is a default handler for not-allowed method
	NotAllowed = http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			writeJSONtoHTTP(errMethodNotAllowed, w, http.StatusNotFound)
		},
	)
)

const defaultContentType = "application/json; charset=utf-8"

// writeHTTPJSON writes the object to the output with the provided http code
// If no object provided, the content of the response would be empty
func writeJSONtoHTTP(obj interface{}, w http.ResponseWriter, c int) {
	w.Header().Set("Content-Type", defaultContentType)
	w.WriteHeader(c)
	if obj == nil {
		return
	}
	if obj, ok := obj.(Error); ok && obj.Code == codeInternal {
		errLog.Println(obj.Reason)
	}
	if err := writeJSON(obj, w); err != nil {
		errLog.Println(err)
	}
}

// writeJSON writes the data to the output buffer
func writeJSON(obj interface{}, w io.Writer) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	if err := enc.Encode(obj); err != nil {
		return fmt.Errorf("write JSON to http: %s", err)
	}
	return nil
}
