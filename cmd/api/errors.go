package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s ", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusInternalServerError, "the server has encountered an  error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s ", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found respobnse: %s path: %s error: %s ", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusNotFound, err.Error())
}
