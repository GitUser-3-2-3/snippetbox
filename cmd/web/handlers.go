package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"snippetbox/pkg/models"
)

func (bknd *backend) home(w http.ResponseWriter, r *http.Request) {
	spt, err := bknd.snippets.Latest()
	if err != nil {
		bknd.serverError(w, err)
		return
	}
	data := bknd.newTemplateData(r)
	data.Snippets = spt
	bknd.renderTemplate(w, r, http.StatusOK, "home.gohtml", data)
}

func (bknd *backend) snippetView(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		bknd.notFound(w)
		return
	}
	spt, err := bknd.snippets.Get(id)
	if errors.Is(err, models.ErrNoRecord) {
		bknd.notFound(w)
		return
	} else if err != nil {
		bknd.serverError(w, err)
		return
	}
	data := bknd.newTemplateData(r)
	data.Snippet = spt
	bknd.renderTemplate(w, r, http.StatusOK, "view.gohtml", data)
}

func (bknd *backend) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := bknd.newTemplateData(r)
	bknd.renderTemplate(w, r, http.StatusOK, "create.gohtml", data)
}

func (bknd *backend) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		bknd.clientError(w, http.StatusBadRequest)
		return
	}
	content := r.PostForm.Get("content")
	title := r.PostForm.Get("title")
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		bknd.clientError(w, http.StatusBadRequest)
		return
	}
	inputErrors := validate(title, content, expires)
	if len(inputErrors) > 0 {
		_, err := fmt.Fprint(w, inputErrors)
		if err != nil {
			bknd.serverError(w, err)
		}
		return
	}
	id, err := bknd.snippets.Insert(title, content, expires)
	if err != nil {
		bknd.serverError(w, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func validate(title, content string, expires int) map[string]string {
	var inputErrors = make(map[string]string)

	if strings.TrimSpace(title) == "" {
		inputErrors["title"] = "This Field cannot be blank"
	} else if utf8.RuneCountInString(title) > 100 {
		inputErrors["title"] = "Title must be less than 100 characters"
	}
	if strings.TrimSpace(content) == "" {
		inputErrors["content"] = "Field cannot be blank"
	}
	if expires != 1 && expires != 7 && expires != 30 && expires != 365 {
		inputErrors["expires"] = "Expires does not match expected value"
	}
	return inputErrors
}
