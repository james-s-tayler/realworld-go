package conduit

import "net/http"

type contextKey string

var userContextKey = contextKey("userContext")

type userContext struct {
	isAuthenticated bool
	userId          int
	username        string
	token           string
}

var anonymousUser = &userContext{}

func (app *Application) getUserContext(r *http.Request) *userContext {
	userContext, ok := r.Context().Value(userContextKey).(*userContext)
	if !ok {
		panic("user context should exist and doesn't")
	}
	return userContext
}
