package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func logInfo(msg string, val ...any) {
	infoLog := log.New(os.Stdout, "INFO::\t", log.Ldate|log.Ltime)
	infoLog.Printf(msg, val)
}

func logError(val ...any) {
	errorLog := log.New(os.Stderr, "ERROR::\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog.Fatal(val)
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address.")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	logInfo("Starting a server on %s", *addr)
	err := http.ListenAndServe(*addr, mux)
	logError(err)
}
