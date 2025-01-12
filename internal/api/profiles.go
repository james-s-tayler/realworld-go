package conduit

import "net/http"

// POST /api/profiles/:username/follow
func (app *Application) followProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// DELETE /api/profiles/:username/follow
func (app *Application) unfollowProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// GET /api/profiles/:username
func (app *Application) getProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}
