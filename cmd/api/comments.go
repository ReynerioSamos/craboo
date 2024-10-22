package main

import (
	_ "encoding/json"
	"fmt"
	"net/http"

	// import the data package which contains the definition for Comment
	// _[space] is used as placeholder to let compiler know to ignore this dependancy
	_ "github.com/ReynerioSamos/craboo/internal/data"
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
	// for now display the result
	fmt.Fprintf(w, "%+v\n", incomingData)
}
