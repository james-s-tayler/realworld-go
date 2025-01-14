package conduit

import "net/http"

// GET /api/articles/:slug
func (app *Application) getArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// GET /api/articles/feed
func (app *Application) getFeedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// POST /api/articles
func (app *Application) createArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// PUT /api/articles/:slug
func (app *Application) updateArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// DELETE /api/articles/:slug
func (app *Application) deleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// POST /api/articles/:slug/comments
func (app *Application) addArticleCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// DELETE /api/articles/:slug/comments/:id
func (app *Application) deleteArticleCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// POST /api/articles/:slug/favorite
func (app *Application) favoriteArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// DELETE /api/articles/:slug/favorite
func (app *Application) unfavoriteArticleHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// GET /api/articles
func (app *Application) getArticlesHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// GET /api/articles/:slug/comments
func (app *Application) getArticleCommentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}
