package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

func (bknd *backend) home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go")

	paths := []string{"./v23/ui/html/base.gohtml",
		"./v23/ui/html/partials/nav.gohtml",
		"./v23/ui/html/pages/home.gohtml",
	}

	tmplt, err := template.ParseFiles(paths...)
	if err != nil {
		bknd.serverError(w, r, err)
		return
	}
	err = tmplt.ExecuteTemplate(w, "base", nil)
	if err != nil {
		bknd.serverError(w, r, err)
	}
}

func (bknd *backend) snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		bknd.notFound(w)
		return
	}
	_, err = fmt.Fprint(w, "Display a snippet with id: ", id)
	if err != nil {
		bknd.serverError(w, r, err)
	}
}

func (bknd *backend) snippetCreate(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Display a form for creating a new snippet"))
	if err != nil {
		bknd.serverError(w, r, err)
	}
}

func (bknd *backend) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte("Creating a new snippet"))
	if err != nil {
		bknd.serverError(w, r, err)
	}
}
