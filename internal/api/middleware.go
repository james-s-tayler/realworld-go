package conduit

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"realworld.tayler.io/internal/data"
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

		tokenString := r.Header.Get("Authorization")
		usercontext := anonymousUser

		if tokenString != "" && strings.HasPrefix(tokenString, "Token ") {
			rawToken := tokenString[len("Token "):]
			token, err := app.tokenService.VerifyToken(rawToken)

			if err != nil {
				app.logger.Error("invalid token", "error", err)
			} else {
				claims, ok := token.Claims.(*data.CustomClaims)
				if ok {
					usercontext = &userContext{
						isAuthenticated: true,
						userId:          claims.UserId,
						username:        claims.Username,
						token:           rawToken,
					}
					// we could validate such a user exists here too etc, but skipping for now
				} else {
					app.logger.Error("there was a problem accessing user claims")
				}
			}
		}

		ctx := context.WithValue(r.Context(), userContextKey, usercontext)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (app *Application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userContext := app.getUserContext(r)

		if !userContext.isAuthenticated {
			app.serveResponseErrorUnauthorized(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
