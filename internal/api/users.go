package conduit

import "net/http"

// GET /api/user
func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// POST /api/users/login
func (app *Application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// POST /api/users
func (app *Application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// PUT /api/user
func (app *Application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}
