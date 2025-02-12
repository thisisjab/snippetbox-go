package main

import (
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /user/signup", dynamic.ThenFunc(app.userSignup))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(app.userSignupPost))
	mux.Handle("GET /user/login", dynamic.ThenFunc(app.userLogin))
	mux.Handle("POST /user/login", dynamic.ThenFunc(app.userLoginPost))

	authRequired := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /snippets/create", authRequired.ThenFunc(app.createSnippet))
	mux.Handle("POST /snippets/create", authRequired.ThenFunc(app.snippetCreatePost))
	mux.Handle("POST /user/logout", authRequired.ThenFunc(app.userLogoutPost))

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippets/view/{id}", dynamic.ThenFunc(app.showSnippet))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
