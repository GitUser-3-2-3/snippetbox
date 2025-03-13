package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (bknd *backend) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = bknd.logError.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (bknd *backend) notFound(w http.ResponseWriter) {
	bknd.clientError(w, http.StatusNotFound)
}

func (bknd *backend) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (bknd *backend) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           bknd.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: bknd.isAuthenticated(r),
	}
}

func (bknd *backend) renderTemplate(w http.ResponseWriter, status int, page string, data *templateData) {
	tmplt, ok := bknd.templateCache[page]
	if !ok {
		bknd.serverError(w, fmt.Errorf("page '%s' not found", page))
		return
	}
	buf := new(bytes.Buffer)
	err := tmplt.ExecuteTemplate(buf, "base", data)
	if err != nil {
		bknd.serverError(w, err)
		return
	}
	w.WriteHeader(status)
	_, err = buf.WriteTo(w)
	if err != nil {
		bknd.serverError(w, fmt.Errorf("error writing content: %w", err))
	}
}

func (bknd *backend) isAuthenticated(r *http.Request) bool {
	return bknd.sessionManager.Exists(r.Context(), "authenticatedUserId")
}
