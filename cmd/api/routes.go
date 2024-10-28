package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (a *applicationDependencies) routes() http.Handler {
	//set up new router
	router := httprouter.New()
	// handle 404
	router.NotFound = http.HandlerFunc(a.notFoundResponse)
	// handle 405
	router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)

	//set up routes
	//route for health checker
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthCheckHandler)

	// routes for comments CRUD functionality
	router.HandlerFunc(http.MethodPost, "/v1/comments", a.createCommentHandler)
	router.HandlerFunc(http.MethodGet, "/v1/comments/:id", a.displayCommentHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/comments/:id", a.updateCommentHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/comments/:id", a.deleteCommentHandler)

	//routes for users CRUD functionality
	router.HandlerFunc(http.MethodPost, "/v1/users", a.createUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", a.displayUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", a.updateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", a.deleteUserHandler)

	//route for List All comments handler
	router.HandlerFunc(http.MethodGet, "/v1/comments", a.ListCommentsHandler)

	//panic recover
	return a.recoverPanic(router)
}
