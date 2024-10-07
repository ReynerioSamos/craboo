package main

import (
	"net/http"
)

func (a *applicationDependencies) healthCheckHandler(w http.ResponseWriter,
	r *http.Request) {

	panic("Apples & Oranges")
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": a.config.environment,
			"version":     appVersion,
		},
	}

	err := a.writeJson(w, http.StatusOK, data, nil)
	if err != nil {
		a.serverErrorResponse(w, r, err)
	}

}
