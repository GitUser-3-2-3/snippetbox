package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
	"snippetbox/internal/validator"
	"snippetbox/pkg/models"
)

type snippetCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"_"`
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
	sptForm := snippetCreateForm{}
	err := bknd.decodePostForm(r, &sptForm)
	if err != nil {
		bknd.clientError(w, http.StatusBadRequest)
		return
	}
	validateForm(&sptForm)
	if !sptForm.Valid() {
		data := bknd.newTemplateData(r)
		data.Form = sptForm
		bknd.renderTemplate(w, http.StatusUnprocessableEntity, "create.gohtml", data)
		return
	}
	id, err := bknd.snippets.Insert(sptForm.Title, sptForm.Content, sptForm.Expires)
	if err != nil {
		bknd.serverError(w, err)
	}
	bknd.sessionManager.Put(r.Context(), "flash", "New snippet created!")
	http.Redirect(w, r, fmt.Sprintf("/snippet/view/%d", id), http.StatusSeeOther)
}

func (bknd *backend) decodePostForm(r *http.Request, dst any) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	err := bknd.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError
		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}
		return err
	}
	return nil
}

func validateForm(sptForm *snippetCreateForm) {
	sptForm.CheckField(validator.NotBlank(sptForm.Title), "title", "Field cannot be blank")

	sptForm.CheckField(validator.MaxChars(sptForm.Title, 100),
		"title", "Field cannot be longer than 100 characters")

	sptForm.CheckField(validator.NotBlank(sptForm.Content), "content", "Field cannot be blank")

	sptForm.CheckField(validator.PermittedInt(
		sptForm.Expires, 1, 30, 365),
		"expires", "Values other than 1, 30, 365 are invalid")
}
