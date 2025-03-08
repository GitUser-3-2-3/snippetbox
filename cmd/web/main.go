package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/form/v4"
	"snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type backend struct {
	templateCache map[string]*template.Template
	logError      *log.Logger
	logInfo       *log.Logger
	snippets      *mysql.SnippetModel
	formDecoder   *form.Decoder
}

func main() {
	dsn := flag.String("dsn", "web:Qwerty1,0*@/snippetbox?parseTime=true", "MySQL database.")
	addr := flag.String("addr", ":4000", "HTTP network address.")
	flag.Parse()

	logInfo := log.New(os.Stdout, "INFO - \t", log.Ldate|log.Ltime)
	logError := log.New(os.Stderr, "ERROR - \t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		logError.Fatal(err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			logError.Fatal("Error closing db", err)
		}
	}(db)
	templateCache, err := newTemplateCache()
	if err != nil {
		logError.Fatal(err)
	}
	formDecoder := form.NewDecoder()
	bknd := &backend{
		templateCache: templateCache,
		logError:      logError,
		logInfo:       logInfo,
		snippets:      &mysql.SnippetModel{DB: db},
		formDecoder:   formDecoder,
	}
	srv := &http.Server{
		ErrorLog:          logError,
		Addr:              *addr,
		Handler:           bknd.routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	logInfo.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	if err != nil {
		logError.Fatal(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
