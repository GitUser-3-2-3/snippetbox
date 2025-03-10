package main

import "net/http"

func (bknd *backend) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		uri    = r.URL.RequestURI()
		method = r.Method
	)
	bknd.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (bknd *backend) notFound(w http.ResponseWriter) {
	bknd.clientError(w, http.StatusNotFound)
}

func (bknd *backend) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
