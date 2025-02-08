package conduit

import (
	"errors"
	"net/http"
	"strconv"

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

	commentId, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		v := validator.New()
		v.AddError("id", "must be an integer")
		app.serveResponseErrorUnprocessableEntity(w, v)
		return
	}

	currentUserId := app.getUserContext(r).userId

	comment, err := app.domains.comments.GetCommentById(int(commentId), currentUserId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrCommentNotFound):
			app.serveResponseErrorNotFound(w, r)
		default:
			app.serveResponseErrorInternalServerError(w, err)
		}
		return
	}

	if comment.ArticleId != article.ArticleId {
		app.serveResponseErrorNotFound(w, r)
		return
	} else if comment.UserId != currentUserId {
		app.serveResponseErrorForbidden(w, r)
		return
	}

	err = app.domains.comments.DeleteComment(comment.CommentId)
	if err != nil {
		app.serveResponseErrorInternalServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GET /api/articles/:slug/comments
func (app *Application) getArticleCommentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello real world"))
}
