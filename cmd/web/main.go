package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

type application struct {
	logError *log.Logger
	logInfo  *log.Logger
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address.")
	flag.Parse()

	logInfo := log.New(os.Stdout, "INFO::\t", log.Ldate|log.Ltime)
	logError := log.New(os.Stderr, "ERROR::\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &application{
		logError: logError,
		logInfo:  logInfo,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: logError,
		Handler:  app.routes(),
	}
	logInfo.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	logError.Fatal(err)
}
