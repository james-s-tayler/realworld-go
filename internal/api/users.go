package conduit

import "net/http"

func (app *Application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}
