package conduit

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	common := alice.New(app.recoverPanic, app.authenticateUser)
	protected := common.Append(app.requireAuthentication)

	// unauthenticated routes
	mux.Handle("POST /api/users/login", common.ThenFunc(app.loginUserHandler))
	mux.Handle("POST /api/users", common.ThenFunc(app.registerUserHandler))
	mux.Handle("GET /api/articles/{slug}", common.ThenFunc(app.getArticleHandler))
	mux.Handle("GET /api/tags", common.ThenFunc(app.getTagsHandler))

	// authenticated routes
	mux.Handle("GET /api/user", protected.ThenFunc(app.getUserHandler))
	mux.Handle("PUT /api/user", protected.ThenFunc(app.updateUserHandler))
	mux.Handle("POST /api/profiles/{username}/follow", protected.ThenFunc(app.followProfileHandler))
	mux.Handle("DELETE /api/profiles/{username}/follow", protected.ThenFunc(app.unfollowProfileHandler))
	mux.Handle("GET /api/articles/feed", protected.ThenFunc(app.getFeedHandler))
	mux.Handle("POST /api/articles", protected.ThenFunc(app.createArticleHandler))
	mux.Handle("PUT /api/articles/{slug}", protected.ThenFunc(app.updateArticleHandler))
	mux.Handle("DELETE /api/articles/{slug}", protected.ThenFunc(app.deleteArticleHandler))
	mux.Handle("POST /api/articles/{slug}/comments", protected.ThenFunc(app.addArticleCommentHandler))
	mux.Handle("DELETE /api/articles/{slug}/comments/{id}", protected.ThenFunc(app.deleteArticleCommentHandler))
	mux.Handle("POST /api/articles/{slug}/favorite", protected.ThenFunc(app.favoriteArticleHandler))
	mux.Handle("DELETE /api/articles/{slug}/favorite", protected.ThenFunc(app.unfavoriteArticleHandler))

	// authentication optional routes
	mux.Handle("GET /api/profiles/{username}", common.ThenFunc(app.getProfileHandler))
	mux.Handle("GET /api/articles", common.ThenFunc(app.getArticlesHandler))
	mux.Handle("GET /api/articles/{slug}/comments", common.ThenFunc(app.getArticleCommentsHandler))

	return mux
}
