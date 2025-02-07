package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"web-dev-journey/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	//files := []string{
	//	"./ui/html/home.page.tmpl",
	//	"./ui/html/base.layout.tmpl",
	//	"./ui/html/footer.partial.tmpl",
	//}
	//
	//ts, err := template.ParseFiles(files...)
	//if err != nil {
	//	app.serverError(w, r, err)
	//}
	//
	//err = ts.Execute(w, nil)
	//if err != nil {
	//	app.serverError(w, r, err)
	//}

	latest, err := strconv.Atoi(r.URL.Query().Get("latest"))

	if err != nil {
		latest = 10
	}

	snippets, err := app.snippets.Latest(latest)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for _, snippet := range snippets {
		fmt.Fprintf(w, "%+v\n", snippet)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))

	if err != nil || id < 1 {
		app.notFound(w)
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
	fmt.Fprintf(w, "%+v", snippet)
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Trying to create a snippet."))
}
