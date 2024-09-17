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
func (a *applicationDependencies) writeJson(w http.ResponseWriter, status int, data any, headers http.Header) error {
	jsResponse, err := json.Marshal(data)
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
	w.Write(jsResponse)

	return nil
}
