package main

import (
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./v23/ui/static/"))
	router.Handle("GET /static/", http.StripPrefix("/static/", fileServer))

	router.HandleFunc("GET /{$}", home)
	router.HandleFunc("GET /snippet/view/{id}", snippetView)
	router.HandleFunc("GET /snippet/create", snippetCreate)
	router.HandleFunc("POST /snippet/create", snippetCreatePost)

	log.Printf("Listening on port :4000")
	log.Fatal(http.ListenAndServe(":4000", router))
}
