package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"snippetbox/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}
	spt, err := app.snippets.Latest()
	if err != nil {
		app.serverError(w, err)
		return
	}
	for _, s := range spt {
		_, _ = fmt.Fprintf(w, "%v\n", s)
	}
	paths := []string{"ui/html/home.page.go.html",
		"ui/html/base.layout.go.html", "ui/html/footer.partial.go.html",
	}
	ts, err := template.ParseFiles(paths...) // 'paths' is a variadic parameter
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, spt)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	spt, err := app.snippets.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		app.notFound(w)
		return
	} else if err != nil {
		app.serverError(w, err)
		return
	}
	data := &templateData{Snippet: spt}

	paths := []string{"./ui/html/show.page.go.html",
		"./ui/html/base.layout.go.html", "./ui/html/footer.partial.go.html",
	}
	ts, err := template.ParseFiles(paths...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	title, content, expires := "O snail", "O snail climb Mount Fuji, but slowly, slowly!", "7"
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
