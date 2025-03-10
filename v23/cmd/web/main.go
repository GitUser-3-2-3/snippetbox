package main

import (
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type backend struct {
	logger *slog.Logger
}

func main() {
	dsn := flag.String("dsn", "web:Qwerty1,0*@/snippetbox?parseTime=true", "MySQL database.")
	addrs := flag.String("addrs", ":4000", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error("Error opening the db", "err", err.Error())
		os.Exit(1)
	}
	if err = run(db, *addrs, logger); err != nil {
		logger.Error("Error starting a server", "err", err.Error())
		if err = db.Close(); err != nil {
			logger.Error("Error closing the db", "err", err.Error())
		}
	}
	os.Exit(1)
}

func run(db *sql.DB, addrs string, logger *slog.Logger) error {
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Error("Error closing the db", "err", err.Error())
		}
	}(db)
	bknd := &backend{logger: logger}
	srvr := http.Server{
		Handler:      bknd.routes(),
		Addr:         addrs,
		ReadTimeout:  5 * time.Second,
		IdleTimeout:  time.Minute,
		WriteTimeout: 10 * time.Second,
	}
	logger.Info("Started a server on", "addrs", addrs)
	return srvr.ListenAndServe()
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	} else if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
