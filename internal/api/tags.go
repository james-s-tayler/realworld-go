package conduit

import "net/http"

// GET /api/tags
func (app *Application) getTagsHandler(w http.ResponseWriter, r *http.Request) {
	tags, err := app.domains.tags.GetAllTags()
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tags": tags}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
	}
}
