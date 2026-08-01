package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/dneedsleep/Social/internal/store"
	"github.com/go-chi/chi/v5"
)

type postKey string

const postContext postKey = "post"

type CreatePostPayload struct {
	Title   string   `json:"title" validate:"required,max=100"`
	Content string   `json:"content" validate:"required,max=1000"`
	Tags    []string `json:"Tags"`
}

type updatePostPayload struct {
	Title   string `json:"title" validate:"required,max=100"`
	Content string `josn:"content" validate:"required,max=1000"`
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {
	var payload CreatePostPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post := &store.Post{
		Title:   payload.Title,
		Content: payload.Content,
		Tags:    payload.Tags,
		// Todo change after auth
		UserID: 1,
	}

	ctx := r.Context()

	if err := app.store.Posts.Create(ctx, post); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.jsonresponse(w, http.StatusCreated, &post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) getPostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	comments, err := app.store.Comments.GetByPostID(r.Context(), post.ID)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	post.Comments = comments

	if err := app.jsonresponse(w, http.StatusCreated, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {
	idParams := chi.URLParam(r, "postID")
	id, err := strconv.ParseInt(idParams, 10, 64)

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	ctx := r.Context()

	err = app.store.Posts.DeleteById(ctx, id)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := app.jsonresponse(w, http.StatusOK, "deleted sucessfully"); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) postContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParams := chi.URLParam(r, "postID")
		id, err := strconv.ParseInt(idParams, 10, 64)

		if err != nil {
			app.internalServerError(w, r, err)
			return
		}

		ctx := r.Context()

		post, err := app.store.Posts.GetById(ctx, id)

		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				WriteJSONError(w, http.StatusNotFound, err.Error())
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, postContext, post)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}

func getPostFromCtx(r *http.Request) *store.Post {
	post, _ := r.Context().Value(postContext).(*store.Post)
	return post
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {
	post := getPostFromCtx(r)

	var payload updatePostPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
	}

	if payload.Content == "" {
		payload.Content = post.Content
	}

	if payload.Title == "" {
		payload.Title = post.Title
	}

	if err := Validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
	}

	post.Content = payload.Content
	post.Title = payload.Title

	if err := app.store.Posts.UpdateById(r.Context(), post); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}

		return
	}

	if err := app.jsonresponse(w, http.StatusOK, post); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
