package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"snippetbox/internal/validator"
	"snippetbox/pkg/models"
)

type snippetCreateForm struct {
	Title   string
	Content string
	Expires int
	validator.Validator
}

func (bknd *backend) home(w http.ResponseWriter, r *http.Request) {
	spt, err := bknd.snippets.Latest()
	if err != nil {
		bknd.serverError(w, err)
		return
	}
	data := bknd.newTemplateData(r)
	data.Snippets = spt
	bknd.renderTemplate(w, http.StatusOK, "home.gohtml", data)
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
	bknd.renderTemplate(w, http.StatusOK, "view.gohtml", data)
}

func (bknd *backend) snippetCreate(w http.ResponseWriter, r *http.Request) {
	data := bknd.newTemplateData(r)
	data.Form = snippetCreateForm{Expires: 365}
	bknd.renderTemplate(w, http.StatusOK, "create.gohtml", data)
}

func (bknd *backend) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		bknd.clientError(w, http.StatusBadRequest)
		return
	}
	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		bknd.clientError(w, http.StatusBadRequest)
		return
	}
	crtForm := snippetCreateForm{
		Title:   r.PostForm.Get("title"),
		Content: r.PostForm.Get("content"),
		Expires: expires,
	}
	validateForm(&crtForm)
	if !crtForm.Valid() {
		data := bknd.newTemplateData(r)
		data.Form = crtForm
		bknd.renderTemplate(w, http.StatusUnprocessableEntity, "create.gohtml", data)
		return
	}
	id, err := bknd.snippets.Insert(crtForm.Title, crtForm.Content, crtForm.Expires)
	if err != nil {
		bknd.serverError(w, err)
	}
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func validateForm(crtForm *snippetCreateForm) {
	crtForm.CheckField(validator.NotBlank(crtForm.Title), "title", "Field cannot be blank")

	crtForm.CheckField(validator.MaxChars(crtForm.Title, 100),
		"title", "Field cannot be longer than 100 characters")

	crtForm.CheckField(validator.NotBlank(crtForm.Content), "content", "Field cannot be blank")

	crtForm.CheckField(validator.PermittedInt(
		crtForm.Expires, 1, 30, 365),
		"expires", "Values other than 1, 30, 365 are invalid")
}
