package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// This Method will accept:
// response write (w)
// status code to send (default is 200)
// actual data to encode in Json
// a map of the headers to set for the response

// create and envelope type
type envelope map[string]any

func (a *applicationDependencies) writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
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
	return nil
}

func (a *applicationDependencies) readJson(w http.ResponseWriter, r *http.Request, destination any) error {
	err := json.NewDecoder(r.Body).Decode(destination)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("the body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("the body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("the body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("the body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("the body must not be empty")
		// the programmer messed up
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		// return other type of error
		default:
			return err
		}
	}
	return nil
}
