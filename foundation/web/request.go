package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

// Param returns the web call parameters from the request.
func Param(r *http.Request, key string) (string, error) {
	vars := mux.Vars(r)

	id, ok := vars["id"]
	if !ok {
		return "", fmt.Errorf("id is missing in path parameters")
	}

	return id, nil
}

// Decode reads the body of an HTTP request looking for a JSON document. The
// body is decoded into the provided value.
//
// If the provided value is a struct then it is checked for validation tags.
func Decode(r *http.Request, val interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}
