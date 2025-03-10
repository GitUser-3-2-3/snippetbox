package main

import (
	"html/template"
	"path/filepath"
	"time"

	"snippetbox/pkg/models"
)

type templateData struct {
	Snippet     *models.Snippet
	CurrentYear int
	Form        any
	Snippets    []*models.Snippet
	Flash       string
}

func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{"humanDate": humanDate}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("./ui/html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)
		tmplt, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.gohtml")
		if err != nil {
			return nil, err
		}
		tmplt, err = tmplt.ParseGlob("./ui/html/partials/*.gohtml")
		if err != nil {
			return nil, err
		}
		tmplt, err = tmplt.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = tmplt
	}
	return cache, nil
}
