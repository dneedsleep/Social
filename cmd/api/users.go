package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/dneedsleep/Social/internal/store"
	"github.com/go-chi/chi/v5"
)

type userKey string

const userContext userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	if err := app.jsonresponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}

}

func (app *application) followUser(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followedID, user.ID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	followerUser := getUserFromContext(r)

	unfollowedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, followerUser.ID, unfollowedID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonresponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParams := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(idParams, 10, 64)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, userID)

		if err != nil {
			switch err {
			case store.ErrNotFound:
				app.notFoundResponse(w, r, store.ErrNotFound)
			default:
				app.internalServerError(w, r, err)
			}
		}

		ctx = context.WithValue(ctx, userContext, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user := r.Context().Value(userContext).(*store.User)
	return user
}

func (app *application) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	pg := store.PaginatedFeedQuery{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	pg, err := pg.Parse(r)
	feed, err := app.store.Posts.GetUserFeed(r.Context(), user, pg)
	if err != nil {
		app.internalServerError(w, r, err)
	}

	if err = app.jsonresponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}

}
