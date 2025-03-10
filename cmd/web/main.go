package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"snippetbox/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type backend struct {
	templateCache  map[string]*template.Template
	logError       *log.Logger
	logInfo        *log.Logger
	snippets       *mysql.SnippetModel
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	dsn := flag.String("dsn", "root:Qwerty1,0*@/snippetbox?parseTime=true", "MySQL database.")
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

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	sessionManager.Store = mysqlstore.New(db)

	bknd := &backend{
		templateCache:  templateCache,
		logError:       logError,
		logInfo:        logInfo,
		snippets:       &mysql.SnippetModel{DB: db},
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
	}
	srv := &http.Server{
		TLSConfig:    tlsConfig,
		ErrorLog:     logError,
		Addr:         *addr,
		Handler:      bknd.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logInfo.Printf("Starting a server on %s", *addr)
	err = srv.ListenAndServeTLS("./ssl/cert.pem", "./ssl/key.pem")
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
