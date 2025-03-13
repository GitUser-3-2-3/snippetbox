package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"snippetbox/pkg/models/mysql"
	"snippetbox/ui"
)

type templateData struct {
	Snippet         *mysql.Snippet
	CurrentYear     int
	Form            any
	Snippets        []*mysql.Snippet
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	local := t.Local()
	return local.Format("02 Jan 2006 at 03:05 PM")
}

var functions = template.FuncMap{"humanDate": humanDate}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.gohtml")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		pathEnd := filepath.Base(page)
		patterns := []string{"html/base.gohtml", "html/partials/*.gohtml", page}
		tmpltSet, err := template.New(pathEnd).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		cache[pathEnd] = tmpltSet
	}
	return cache, nil
}
