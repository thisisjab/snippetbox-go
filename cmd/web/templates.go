package main

import (
	"html/template"
	"path/filepath"
	"time"
	"web-dev-journey/internal/model"
)

type templateData struct {
	CurrentYear     int
	Snippet         model.Snippet
	Snippets        []model.Snippet
	Flash           string
	Form            any
	IsAuthenticated bool
}

func humanDateTime(t time.Time) string {
	return t.Format("2006-01-02 3:04 PM")
}

var funcMap = template.FuncMap{
	"humanDateTime": humanDateTime,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := filepath.Glob("./ui/html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(funcMap).ParseFiles("./ui/html/base.gohtml")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.gohtml")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
