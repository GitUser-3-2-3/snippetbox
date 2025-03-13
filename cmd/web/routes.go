package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (bknd *backend) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		bknd.notFound(w)
	})
	server := http.FileServer(http.Dir("./ui/static"))
	router.Handler(http.MethodGet, "/static/*path", http.StripPrefix("/static", server))

	dynamic := alice.New(bknd.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(bknd.home))
	router.Handler(http.MethodGet, "/snippet/view/:id", dynamic.ThenFunc(bknd.snippetView))
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(bknd.userSignUp))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(bknd.userSignUpPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(bknd.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(bknd.userLoginPost))

	protected := dynamic.Append(bknd.requireAuthentication)

	router.Handler(http.MethodGet, "/snippet/create", protected.ThenFunc(bknd.snippetCreate))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(bknd.userLogoutPost))
	router.Handler(http.MethodPost, "/snippet/create", protected.ThenFunc(bknd.snippetCreatePost))

	standard := alice.New(bknd.recoverPanic, bknd.logRequest, secureHeaders)
	return standard.Then(router)
}
