package conduit

import (
	"encoding/json"
	"fmt"
	"net/http"

	"realworld.tayler.io/internal/data"
)

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

	var input struct {
		User struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"user"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		msg := fmt.Sprintf("An unexpected error occurred while processing the request: %v\n", err.Error())
		app.logger.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}

	user := &data.User{
		Username: input.User.Username,
		Email:    input.User.Email,
		Token:    input.User.Password,
		Bio:      "I work at statefarm",
		Image:    nil,
	}

	// TODO: validate input

	user, err = app.domains.users.RegisterUser(user)
	if err != nil {
		msg := fmt.Sprintf("An unexpected error occurred while processing the request: %v\n", err.Error())
		app.logger.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}

	app.logger.Info("Registered user: ", "user", user)

	err = json.NewEncoder(w).Encode(envelope{"user": user})
	if err != nil {
		msg := fmt.Sprintf("An unexpected error occurred while processing the request: %v\n", err.Error())
		app.logger.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
}

// PUT /api/user
func (app *Application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}
