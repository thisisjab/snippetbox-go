package main

import (
	"errors"
	"html/template"
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

	files := []string{
		"./ui/html/pages/home.tmpl",
		"./ui/html/base.tmpl",
		"./ui/html/partials/footer.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
	}

	err = ts.Execute(w, templateData{Snippets: snippets})
	if err != nil {
		app.serverError(w, r, err)
	}
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

	files := []string{
		"./ui/html/base.tmpl",
		"./ui/html/partials/nav.tmpl",
		"./ui/html/pages/view.tmpl",
		"./ui/html/partials/footer.tmpl",
	}

	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = ts.ExecuteTemplate(w, "base", templateData{Snippet: snippet})
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Trying to create a snippet."))
}
