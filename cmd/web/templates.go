package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"
	"web-dev-journey/internal/model"
	"web-dev-journey/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         model.Snippet
	Snippets        []model.Snippet
	Flash           string
	Form            any
	IsAuthenticated bool
	CSRFToken       string
}

func humanDateTime(t time.Time) string {
	return t.Format("2006-01-02 3:04 PM")
}

var funcMap = template.FuncMap{
	"humanDateTime": humanDateTime,
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.Files, "html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.gohtml",
			"html/partials/*.gohtml",
			page,
		}

		ts, err := template.New(name).Funcs(funcMap).ParseFS(ui.Files, patterns...)

		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}
