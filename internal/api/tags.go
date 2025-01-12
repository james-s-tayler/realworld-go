package conduit

import "net/http"

// GET /api/tags
func (app *Application) getTagsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}
