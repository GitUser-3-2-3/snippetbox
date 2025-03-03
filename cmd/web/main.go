package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"

	"snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logError      *log.Logger
	logInfo       *log.Logger
	snippets      *mysql.SnippetModel
	templateCache map[string]*template.Template
}

func main() {
	dsn := flag.String("dsn", "web:Qwerty1,0*@/snippetbox?parseTime=true", "MySQL database.")
	addr := flag.String("addr", ":4000", "HTTP network address.")
	flag.Parse()

	logInfo := log.New(os.Stdout, "INFO::\t", log.Ldate|log.Ltime)
	logError := log.New(os.Stderr, "ERROR::\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		logError.Fatal(err)
	}
	defer func(db *sql.DB) {
		if err := db.Close(); err != nil {
			logError.Fatal("Error closing db", err)
		}
	}(db)
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		logError.Fatal(err)
	}
	app := &application{
		logError:      logError,
		logInfo:       logInfo,
		snippets:      &mysql.SnippetModel{DB: db},
		templateCache: templateCache,
	}
	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: logError,
		Handler:  app.routes(),
	}
	logInfo.Printf("Starting server on %s", *addr)
	err = srv.ListenAndServe()
	logError.Fatal(err)
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
