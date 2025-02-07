package main

import (
	"errors"
	"net/http"
	"strconv"
	"web-dev-journey/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	latest, err := strconv.Atoi(r.URL.Query().Get("latest"))

	if err != nil {
		latest = 10
	}

	snippets, err := app.snippets.Latest(latest)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.render(w, r, http.StatusOK, "home.tmpl", templateData{Snippets: snippets})
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.render(w, r, http.StatusOK, "view.tmpl", templateData{Snippet: snippet})
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Trying to create a snippet."))
}
