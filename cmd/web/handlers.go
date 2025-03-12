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

type userSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"_"`
}

func (bknd *backend) userSignUp(w http.ResponseWriter, r *http.Request) {
	data := bknd.newTemplateData(r)
	data.Form = userSignUpForm{}
	bknd.renderTemplate(w, http.StatusOK, "signup.gohtml", data)
}

func (bknd *backend) userSignUpPost(w http.ResponseWriter, r *http.Request) {
	signUpForm := userSignUpForm{}
	err := bknd.decodePostForm(r, &signUpForm)
	if err != nil {
		bknd.clientError(w, http.StatusBadRequest)
		return
	}
	validateSignUpForm(&signUpForm)
	if !signUpForm.Valid() {
		data := bknd.newTemplateData(r)
		data.Form = signUpForm
		bknd.renderTemplate(w, http.StatusUnprocessableEntity, "signup.gohtml", data)
		return
	}
	err = bknd.users.Insert(signUpForm.Name, signUpForm.Email, signUpForm.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			signUpForm.AddFieldError("email", "Email already in use!")
			data := bknd.newTemplateData(r)
			data.Form = signUpForm
			bknd.renderTemplate(w, http.StatusUnprocessableEntity, "signup.gohtml", data)
		} else {
			bknd.serverError(w, err)
		}
		return
	}
	bknd.sessionManager.Put(r.Context(), "flash", "New user signed up! Please log in!")
	http.Redirect(w, r, "/user/login/", http.StatusSeeOther)
}

func validateSignUpForm(signupForm *userSignUpForm) {
	signupForm.CheckField(validator.NotBlank(signupForm.Name), "name", "Field cannot be blank")

	signupForm.CheckField(validator.NotBlank(signupForm.Email), "email", "Field cannot be blank")
	signupForm.CheckField(validator.Matches(signupForm.Email, validator.EmailRX),
		"email", "Field must be a valid email address")

	signupForm.CheckField(validator.NotBlank(signupForm.Password), "password", "Field cannot be blank")
	signupForm.CheckField(validator.MinChars(signupForm.Password, 8),
		"password", "Password must be at least 8 characters long")
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"_"`
}

func (bknd *backend) userLogin(w http.ResponseWriter, r *http.Request) {
	data := bknd.newTemplateData(r)
	data.Form = userLoginForm{}
	bknd.renderTemplate(w, http.StatusOK, "login.gohtml", data)
}

func (bknd *backend) userLoginPost(w http.ResponseWriter, r *http.Request) {
	loginForm := userLoginForm{}
	err := bknd.decodePostForm(r, &loginForm)
	if err != nil {
		bknd.clientError(w, http.StatusBadRequest)
		return
	}
	validateLoginForm(&loginForm)
	if !loginForm.Valid() {
		data := bknd.newTemplateData(r)
		data.Form = loginForm
		bknd.renderTemplate(w, http.StatusUnprocessableEntity, "login.gohtml", data)
		return
	}
	id, err := bknd.users.Authenticate(loginForm.Email, loginForm.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			loginForm.AddNonFieldError("Email or Password is invalid!")
			data := bknd.newTemplateData(r)
			data.Form = loginForm
			bknd.renderTemplate(w, http.StatusUnprocessableEntity, "login.gohtml", data)
		} else {
			bknd.serverError(w, err)
		}
		return
	}
	err = bknd.sessionManager.RenewToken(r.Context())
	if err != nil {
		bknd.serverError(w, err)
		return
	}
	bknd.sessionManager.Put(r.Context(), "authenticatedUserId", id)
	http.Redirect(w, r, "/snippet/create", http.StatusSeeOther)
}

func validateLoginForm(loginForm *userLoginForm) {
	loginForm.CheckField(validator.NotBlank(loginForm.Email), "email", "Field cannot be blank")
	loginForm.CheckField(validator.Matches(loginForm.Email, validator.EmailRX),
		"email", "Field must be a valid email address")

	loginForm.CheckField(validator.NotBlank(loginForm.Password), "password", "Field cannot be blank")
}

func (bknd *backend) userLogoutPost(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprintln(w, "Logout")
}
