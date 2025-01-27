package conduit

import (
	"errors"
	"net/http"

	"realworld.tayler.io/internal/data"
)

// POST /api/profiles/:username/follow
func (app *Application) followProfileHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")

	lookupUser, err := app.domains.users.GetUserByUsername(username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUserNotFound):
			app.serveResponseErrorNotFound(w, r)
		default:
			app.serveResponseErrorInternalServerError(w, err)
		}
		return
	}

	err = app.domains.users.Follow(app.getUserContext(r).userId, lookupUser.UserId)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	profile := &data.Profile{
		Username:  lookupUser.Username,
		Bio:       lookupUser.Bio,
		Image:     lookupUser.Image,
		Following: true,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
	}
}

// DELETE /api/profiles/:username/follow
func (app *Application) unfollowProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// GET /api/profiles/:username
func (app *Application) getProfileHandler(w http.ResponseWriter, r *http.Request) {

	username := r.PathValue("username")

	lookupUser, err := app.domains.users.GetUserByUsername(username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUserNotFound):
			app.serveResponseErrorNotFound(w, r)
		default:
			app.serveResponseErrorInternalServerError(w, err)
		}
		return
	}

	profile := &data.Profile{
		Username:  lookupUser.Username,
		Bio:       lookupUser.Bio,
		Image:     lookupUser.Image,
		Following: false,
	}

	if userContext := app.getUserContext(r); userContext.isAuthenticated {
		isFollowing, err := app.domains.users.IsFollowing(userContext.userId, lookupUser.UserId)
		if err != nil {
			app.serveResponseErrorInternalServerError(w, err)
			return
		}
		profile.Following = isFollowing
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"profile": profile}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
	}
}
