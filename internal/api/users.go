package conduit

import (
	"encoding/json"
	"errors"
	"net/http"

	"realworld.tayler.io/internal/data"
	"realworld.tayler.io/internal/validator"
)

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
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	user := &data.User{
		Username: input.User.Username,
		Email:    input.User.Email,
		Bio:      "I work at statefarm",
		Image:    nil,
	}

	err = user.Password.Set(input.User.Password)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	v := validator.New()
	if user.Validate(v); !v.Valid() {
		app.serveResponseErrorUnprocessableEntity(w, v)
		return
	}

	user, err = app.domains.users.RegisterUser(user)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	token, err := app.tokenService.CreateToken(user)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}
	user.Token = token

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}
}

// POST /api/users/login
func (app *Application) loginUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		User struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"user"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	user, err := app.domains.users.GetUserByCredentials(input.User.Email, input.User.Password)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrInvalidCredentials):
			app.serveResponseErrorUnauthorized(w, r)
			return
		default:
			app.serveResponseErrorInternalServerError(w, err)
			return
		}
	}

	token, err := app.tokenService.CreateToken(user)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}
	user.Token = token

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
	}
}

// GET /api/user
func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {

	userContext := app.getUserContext(r)
	user, err := app.domains.users.GetUserById(userContext.userId)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}
	user.Token = userContext.token

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}
}

// PUT /api/user
func (app *Application) updateUserHandler(w http.ResponseWriter, r *http.Request) {

	userContext := app.getUserContext(r)
	user, err := app.domains.users.GetUserById(userContext.userId)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}
	user.Token = userContext.token

	var input struct {
		User struct {
			Username *string `json:"username"`
			Email    *string `json:"email"`
			Password *string `json:"password"`
			Image    *string `json:"image"`
			Bio      *string `json:"bio"`
		} `json:"user"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	if input.User.Username != nil {
		user.Username = *input.User.Username
	}
	if input.User.Email != nil {
		user.Email = *input.User.Email
	}
	if input.User.Password != nil {
		user.Password.Set(*input.User.Password)
	}
	if input.User.Image != nil {
		user.Image = input.User.Image
	}
	if input.User.Bio != nil {
		user.Bio = *input.User.Bio
	}

	v := validator.New()
	if user.Validate(v); !v.Valid() {
		app.serveResponseErrorUnprocessableEntity(w, v)
		return
	}

	err = app.domains.users.UpdateUser(user)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}
}
