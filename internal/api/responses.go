package conduit

import (
	"fmt"
	"net/http"
)

func (app *Application) serveResponseErrorInternalServerError(w http.ResponseWriter, err error) {
	msg := fmt.Sprintf("An unexpected error occurred while processing the request: %v\n", err.Error())
	app.logger.Error(msg)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, msg)
}

func (app *Application) serveResponseErrorUnauthorized(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Unauthorized request to %v %v from ip address: %v\n", r.Method, r.RequestURI, r.RemoteAddr)
	app.logger.Warn(msg)
	w.WriteHeader(http.StatusUnauthorized)
}
