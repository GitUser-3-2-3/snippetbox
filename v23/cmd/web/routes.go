package main

import (
	"net/http"
)

func (bknd *backend) routes() http.Handler {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./v23/ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	router.HandleFunc("GET /{$}", bknd.home)
	router.HandleFunc("GET /snippet/view/{id}", bknd.snippetView)
	router.HandleFunc("GET /snippet/create", bknd.snippetCreate)
	router.HandleFunc("POST /snippet/create", bknd.snippetCreatePost)

	return router
}
