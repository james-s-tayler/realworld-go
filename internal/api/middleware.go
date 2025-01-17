package conduit

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
)

// panic recovery middleware
// context middleware
// require auth middleware

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				error := fmt.Errorf("%s", err)
				app.logger.Error(error.Error(), slog.String("method", r.Method), slog.String("uri", r.RequestURI))

				w.Header().Set("Connection", "close")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := context.WithValue(r.Context(), userContextKey, anonymousUser)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (app *Application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userContext := r.Context().Value(userContextKey).(*userContext)

		if !userContext.isAuthenticated {
			app.serveResponseErrorUnauthorized(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
