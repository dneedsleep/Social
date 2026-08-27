package main

import (
	"net/http"
)

func (app *application) healthCheckerHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonresponse(w, http.StatusOK, data); err != nil {
		WriteJSONError(w, http.StatusInternalServerError, err.Error())
	}
}
