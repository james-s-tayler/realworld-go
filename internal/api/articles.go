package conduit

import (
	"errors"
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
	slug := r.PathValue("slug")

	article, err := app.domains.articles.GetArticleBySlug(slug, app.getUserContext(r).userId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrArticleNotFound):
			app.serveResponseErrorNotFound(w, r)
		default:
			app.serveResponseErrorInternalServerError(w, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"article": article}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
	}
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
		switch {
		case errors.Is(err, data.ErrDuplicateSlug):
			v.AddError("slug", "duplicate slug")
			app.serveResponseErrorUnprocessableEntity(w, v)
		default:
			app.serveResponseErrorInternalServerError(w, err)
		}
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

	slug := r.PathValue("slug")

	article, err := app.domains.articles.GetArticleBySlug(slug, app.getUserContext(r).userId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrArticleNotFound):
			app.serveResponseErrorNotFound(w, r)
		default:
			app.serveResponseErrorInternalServerError(w, err)
		}
		return
	}

	currentUserId := app.getUserContext(r).userId
	if article.UserId != currentUserId {
		app.serveResponseErrorForbidden(w, r)
		return
	}

	err = app.domains.articles.DeleteArticle(article.ArticleId, currentUserId)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// POST /api/articles/:slug/favorite
func (app *Application) favoriteArticleHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	article, err := app.domains.articles.GetArticleBySlug(slug, app.getUserContext(r).userId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrArticleNotFound):
			app.serveResponseErrorNotFound(w, r)
		default:
			app.serveResponseErrorInternalServerError(w, err)
		}
		return
	}

	if !article.Favorited {
		err = app.domains.articles.FavoriteArticle(article.ArticleId, app.getUserContext(r).userId)
		if err != nil {
			app.serveResponseErrorInternalServerError(w, err)
			return
		}

		article, err = app.domains.articles.GetArticleBySlug(slug, app.getUserContext(r).userId)
		if err != nil {
			app.serveResponseErrorInternalServerError(w, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"article": article}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
	}
}

// DELETE /api/articles/:slug/favorite
func (app *Application) unfavoriteArticleHandler(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	article, err := app.domains.articles.GetArticleBySlug(slug, app.getUserContext(r).userId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrArticleNotFound):
			app.serveResponseErrorNotFound(w, r)
		default:
			app.serveResponseErrorInternalServerError(w, err)
		}
		return
	}

	if article.Favorited {
		err = app.domains.articles.UnfavoriteArticle(article.ArticleId, app.getUserContext(r).userId)
		if err != nil {
			app.serveResponseErrorInternalServerError(w, err)
			return
		}

		article, err = app.domains.articles.GetArticleBySlug(slug, app.getUserContext(r).userId)
		if err != nil {
			app.serveResponseErrorInternalServerError(w, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"article": article}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
	}
}
