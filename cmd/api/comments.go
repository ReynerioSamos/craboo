package main

import (
	_ "encoding/json"
	"errors"
	"fmt"
	"net/http"

	// import the data package which contains the definition for Comment
	// _[space] is used as placeholder to let compiler know to ignore this dependancy
	"github.com/ReynerioSamos/craboo/internal/data"
	"github.com/ReynerioSamos/craboo/internal/validator"
)

func (a *applicationDependencies) createCommentHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold a comment
	// we use struct tags [``] to make the names display in lowercase
	var incomingData struct {
		Content string `json:"content"`
		Author  string `json:"author"`
	}

	// perform the decoding
	err := a.readJson(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Copy the values from incomingData to a new Comment struct
	// At this point in our code the JSON is well-formed JSON so now
	// we will validate it using the Validators which expects a Comment
	comment := &data.Comment{
		Content: incomingData.Content,
		Author:  incomingData.Author,
	}
	// Intialize Validator instance
	v := validator.New()
	// Do the validation
	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Add the comment to the database table
	err = a.commentModel.Insert(comment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// for now display the result
	fmt.Fprintf(w, "%+v\n", incomingData)

	// Set a Location header. The path to the newly created comment
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/comments/%d", comment.ID))

	// Send a JSON response with 201 (new resource created) status code
	data := envelope{
		"comment": comment,
	}
	err = a.writeJson(w, http.StatusCreated, data, headers)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) displayCommentHandler(w http.ResponseWriter, r *http.Request) {
	// get the id from the URL /v1/comments/:id so that we can use it to query the comments table
	// We will implement the readIDParam() function later
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	//Call Get() to retrieve the comment with the specified id
	comment, err := a.commentModel.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	// display the comment
	data := envelope{
		"comment": comment,
	}
	err = a.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) updateCommentHandler(w http.ResponseWriter, r *http.Request) {
	// Get the id from the URL
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	// Call Get() to retrieve the comment with specified id
	comment, err := a.commentModel.Get(id)
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
		Content *string `json:"content"`
		Author  *string `json:"author"`
	}

	// decoding
	err = a.readJson(w, r, &incomingData)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// We need to now check the fields to see which ones need updating
	// if incomingData.Content is nil, no update done
	if incomingData.Content != nil {
		comment.Content = *incomingData.Content
	}

	// if incomingData.Author is nil, no updatr was provided
	if incomingData.Author != nil {
		comment.Author = *incomingData.Author
	}

	// Before we write the updates to the DB let's validate
	v := validator.New()
	data.ValidateComment(v, comment)
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// perform the update
	err = a.commentModel.Update(comment)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	data := envelope{
		"comment": comment,
	}
	err = a.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
}

func (a *applicationDependencies) deleteCommentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := a.readIDParam(r)
	if err != nil {
		a.notFoundResponse(w, r)
		return
	}

	err = a.commentModel.Delete(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			a.notFoundResponse(w, r)
		default:
			a.serverErrorResponse(w, r, err)
		}
		return
	}
	// display the comment
	data := envelope{
		"message": "comment successfully deleted",
	}
	err = a.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}
}

// create the list handler
func (a *applicationDependencies) ListCommentsHandler(w http.ResponseWriter, r *http.Request) {
	// create a struct to hold the query parameters
	// later we will add fields for pagination and sorting (filters)
	var queryParametersData struct {
		Content string
		Author  string
	}
	// get the query parameters from the URL
	queryParameters := r.URL.Query()

	// load the query parameters into our struct
	queryParametersData.Content = a.getSingleQueryParameter(
		queryParameters,
		"content", "")

	queryParametersData.Author = a.getSingleQueryParameter(
		queryParameters,
		"author", "")

	comments, err := a.commentModel.GetAll(queryParametersData.Content, queryParametersData.Author)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	data := envelope{
		"comments": comments,
	}

	err = a.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}

}
