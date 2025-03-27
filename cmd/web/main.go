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

	"snippetbox/pkg/models"
	"snippetbox/pkg/models/mysql"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"

	_ "github.com/go-sql-driver/mysql"
)

type backend struct {
	templateCache  map[string]*template.Template
	logError       *log.Logger
	logInfo        *log.Logger
	snippets       *mysql.SnippetModel
	users          *models.UserModel
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
}

func main() {
	dsn := flag.String("dsn", "root:Qwerty1,0*@/snippetbox?parseTime=true", "MySQL database.")
	addr := flag.String("addr", ":4000", "HTTP network address.")
	flag.Parse()

	logInfo := log.New(os.Stdout, "INFO - \t", log.Ldate|log.Ltime)
	logError := log.New(os.Stderr, "ERROR - \t", log.Ldate|log.Ltime|log.Lshortfile)

	bknd, err := backendInit(*dsn, logInfo, logError)
	if err != nil {
		logError.Fatal(err)
	}
	defer bknd.closeDB()

	err = runServer(*addr, bknd, logInfo, logError)
	if err != nil {
		logError.Fatal(err)
	}
}

func backendInit(dsn string, logInfo *log.Logger, logError *log.Logger) (*backend, error) {
	db, err := openDB(dsn)
	if err != nil {
		return nil, err
	}
	templateCache, err := newTemplateCache()
	if err != nil {
		return nil, err
	}
	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true
	sessionManager.Store = mysqlstore.New(db)

	return &backend{
		templateCache:  templateCache,
		logError:       logError,
		logInfo:        logInfo,
		snippets:       &mysql.SnippetModel{DB: db},
		users:          &models.UserModel{DB: db},
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}, nil
}

func runServer(addrs string, bknd *backend, logInfo *log.Logger, logError *log.Logger) error {
	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:       tls.VersionTLS12,
	}
	srv := &http.Server{
		TLSConfig:    tlsConfig,
		ErrorLog:     logError,
		Addr:         addrs,
		Handler:      bknd.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logInfo.Printf("Starting a server on %s", addrs)
	return srv.ListenAndServeTLS("./ssl/cert.pem", "./ssl/key.pem")
}

func (bknd *backend) closeDB() {
	if bknd.snippets != nil && bknd.snippets.DB != nil {
		if err := bknd.snippets.DB.Close(); err != nil {
			bknd.logError.Fatal("Error closing db", err)
		}
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
