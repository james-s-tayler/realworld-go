package conduit

type contextKey string

var userContextKey = contextKey("userContext")

type userContext struct {
	isAuthenticated bool
	userId          int
	username        string
}

var anonymousUser = &userContext{}
