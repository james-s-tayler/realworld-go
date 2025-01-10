package conduit

import "net/http"

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/user", app.getUserHandler)
	return mux
}
