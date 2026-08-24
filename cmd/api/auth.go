package main

import (
	"net/http"

	"github.com/dneedsleep/Social/internal/store"
)

type RegisteredUserPayload struct {
	FirstName string `json:"first_name" validate:"required,max=100"`
	LastName  string `json:"last_name" validate:"required,max=100"`
	Email     string `json:"email" validate:"required,max=255"`
	Password  string `json:"password" validate:"required,max=72"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisteredUserPayload
	if err := readJSON(w, r, payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)

		return
	}

	user := &store.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
	}

	// hash the password

	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonresponse(w, http.StatusCreated, nil); err != nil {
		app.internalServerError(w, r, err)

	}

}
