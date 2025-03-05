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

	router.HandlerFunc(http.MethodGet, "/", bknd.home)
	router.HandlerFunc(http.MethodGet, "/snippet/view/:id", bknd.snippetView)
	router.HandlerFunc(http.MethodGet, "/snippet/create", bknd.snippetCreate)
	router.HandlerFunc(http.MethodPost, "/snippet/create", bknd.snippetCreatePost)

	standard := alice.New(bknd.recoverPanic, bknd.logRequest, secureHeaders)
	return standard.Then(router)
}
