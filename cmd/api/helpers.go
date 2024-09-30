package main

import (
	"encoding/json"
	"net/http"
)

// This Method will accept:
// response write (w)
// status code to send (default is 200)
// actual data to encode in Json
// a map of the headers to set for the response
<<<<<<< HEAD
func (a *applicationDependencies) writeJson(w http.ResponseWriter, status int, data any, headers http.Header) error {
=======

// create and envelope type
type envelope map[string]any

func (a *applicationDependencies) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
>>>>>>> encode-json-v5
	jsResponse, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	jsResponse = append(jsResponse, '\n')
	for key, value := range headers {
		w.Header()[key] = value
		//w.Header().Set(key, value[0])
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jsResponse)
	if err != nil {
		return err
	}
<<<<<<< HEAD
=======

>>>>>>> encode-json-v5
	return nil
}
