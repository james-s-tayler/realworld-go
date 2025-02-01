package conduit

import (
	"net/http"

	"realworld.tayler.io/internal/data"
	"realworld.tayler.io/internal/validator"
)

// GET /api/articles
func (app *Application) getArticlesHandler(w http.ResponseWriter, r *http.Request) {
	filters := &data.ArticleFilters{}
	filters.ParseFilters(r)

	articles := make([]data.BodylessArticle, 0)

	err := app.writeJSON(w, http.StatusOK, envelope{"articles": articles, "articlesCount": len(articles)}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
	}
}

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
	var input data.CreateArticleDTO

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	v := validator.New()

	if input.Validate(v); !v.Valid() {
		app.serveResponseErrorUnprocessableEntity(w, v)
		return
	}

	article, err := app.domains.articles.CreateArticle(input, app.getUserContext(r).userId)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"article": article}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
	}
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

// GET /api/articles/:slug/comments
func (app *Application) getArticleCommentsHandler(w http.ResponseWriter, r *http.Request) {
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
