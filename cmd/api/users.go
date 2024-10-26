package main

import (
	_ "encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ReynerioSamos/craboo/internal/data"
	"github.com/ReynerioSamos/craboo/internal/validator"
)

func (a *applicationDependencies) createUserHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold a User
	var incomingData struct {
		Email    string `json:"email"`
		Fullname string `json:"fullname"`
	}

	// decoding
	err := a.readJson(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from incomingData to a new User struct
	// we will validate it using the Validators which expects a User
	user := &data.User{
		Email:    incomingData.Email,
		Fullname: incomingData.Fullname,
	}
	// Intialize Validator instance
	v := validator.New()
	// Do the validation
	data.ValidateUser(v, user)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Add the user to the database table
	err = a.userModel.Insert(user)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// for now display the result
	fmt.Fprintf(w, "%+v\n", incomingData)

	// Set a Location header. The path to the newly created user
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/users/%d", user.ID))

	// Send a JSON response with 201 (new resource created) status code
	data := envelope{
		"user": user,
	}
	err = a.writeJson(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) displayUserHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the URL /v1/users/:id so that we can use it to query the users table
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	//Call Get() to retrieve the User with the specified id
	user, err := a.userModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	// display the User
	data := envelope{
		"user": user,
	}
	err = a.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Call Get() to retrieve the User with specified id
	user, err := a.userModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}

	// Use our temporary incomingData struct to hold the data
	// Note: types have been changed to pointers to differentiate b/w the client
	// leaving a field empty intentionally and the field not needing to be updated

	var incomingData struct {
		Email    *string `json:"email"`
		Fullname *string `json:"fullname"`
	}

	// decoding
	err = a.readJson(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// We need to now check the fields to see which ones need updating
	// if incomingData.Email is nil, no update done
	if incomingData.Email != nil {
		user.Email = *incomingData.Email
	}

	// if incomingData.Fullname is nil, no update was provided
	if incomingData.Fullname != nil {
		user.Fullname = *incomingData.Fullname
	}

	// Before we write the updates to the DB let's validate
	v := validator.New()
	data.ValidateUser(v, user)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// perform the update
	err = a.userModel.Update(user)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	data := envelope{
		"user": user,
	}
	err = a.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.userModel.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	// display the user
	data := envelope{
		"message": "user successfully deleted",
	}
	err = a.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}
