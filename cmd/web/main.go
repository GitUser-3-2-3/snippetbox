package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	logInfo("Starting a server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func logInfo(log string) {
	fmt.Printf("INFO:: %s", log)
}
