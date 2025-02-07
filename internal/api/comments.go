package conduit

import (
	"errors"
	"net/http"

	"realworld.tayler.io/internal/data"
	"realworld.tayler.io/internal/validator"
)

// POST /api/articles/:slug/comments
func (app *Application) addArticleCommentHandler(w http.ResponseWriter, r *http.Request) {
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

	var input data.CreateCommentDTO

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	v := validator.New()

	if input.Validate(v); !v.Valid() {
		app.serveResponseErrorUnprocessableEntity(w, v)
		return
	}

	userId := app.getUserContext(r).userId

	comment, err := app.domains.comments.CreateComment(article.ArticleId, userId, *input.Comment.Body)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"comment": comment}, nil)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}
}

// DELETE /api/articles/:slug/comments/:id
func (app *Application) deleteArticleCommentHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}

// GET /api/articles/:slug/comments
func (app *Application) getArticleCommentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}
