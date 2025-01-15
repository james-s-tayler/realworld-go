package conduit

import (
	"encoding/json"
	"fmt"
	"net/http"

	"realworld.tayler.io/internal/data"
)

// GET /api/user
func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := &data.User{}
	err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		msg := fmt.Sprintf("An unexpected error occurred while processing the request: %v\n", err.Error())
		app.logger.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
}

// POST /api/users/login
func (app *Application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	user := &data.User{}
	err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		msg := fmt.Sprintf("An unexpected error occurred while processing the request: %v\n", err.Error())
		app.logger.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
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
		Bio:      "I work at statefarm",
		Image:    nil,
	}

	// TODO: validate input

	err = user.Password.Set(input.User.Password)
	if err != nil {
		msg := fmt.Sprintf("An unexpected error occurred while processing the request: %v\n", err.Error())
		app.logger.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}

	user, err = app.domains.users.RegisterUser(user)
	if err != nil {
		msg := fmt.Sprintf("An unexpected error occurred while processing the request: %v\n", err.Error())
		app.logger.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}

	token, err := app.tokenService.CreateToken(user)
	if err != nil {
		msg := fmt.Sprintf("An unexpected error occurred while processing the request: %v\n", err.Error())
		app.logger.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
	user.Token = token

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
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
	user := &data.User{}
	err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		msg := fmt.Sprintf("An unexpected error occurred while processing the request: %v\n", err.Error())
		app.logger.Error(msg)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, msg)
		return
	}
}
