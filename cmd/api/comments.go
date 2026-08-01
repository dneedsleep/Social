package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dneedsleep/Social/internal/store"
	"github.com/go-chi/chi/v5"
)

type CreateCommentPayload struct {
	Content string `json:"content" validate:"required,max=1000"`
}

func (app *application) CreateCommentHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreateCommentPayload

	uID := r.Header.Get("X-User-Id")
	if uID == "" {
		app.badRequestResponse(w, r, errors.New("missing X-User-Id header"))
		return
	}

	userId, err := strconv.ParseInt(uID, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid X-User-Id"))
		return
	}

	pID := chi.URLParam(r, "postID")
	postId, err := strconv.ParseInt(pID, 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, errors.New("invalid postID"))
		return
	}

	if err := readJSON(w, r, &payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	comment := &store.Comment{
		PostID:  postId,
		UserID:  userId,
		Content: payload.Content,
	}

	ctx := r.Context()

	if err := app.store.Comments.Create(ctx, comment); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := app.jsonresponse(w, http.StatusCreated, comment); err != nil {
		app.internalServerError(w, r, err)
	}

}
