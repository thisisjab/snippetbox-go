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

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippets/view/{id}", dynamic.ThenFunc(app.showSnippet))
	mux.Handle("GET /snippets/create", dynamic.ThenFunc(app.createSnippet))
	mux.Handle("POST /snippets/create", dynamic.ThenFunc(app.snippetCreatePost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
