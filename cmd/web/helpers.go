package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	_ = app.logError.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), 500)
}

func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

//goland:noinspection GoUnusedParameter
func (app *application) render(
	w http.ResponseWriter, r *http.Request, name string, td *templateData) {

	tmplt, ok := app.templateCache[name]
	if !ok {
		app.serverError(
			w, fmt.Errorf("the template %s does not exist", name))
		return
	}
	buf := new(bytes.Buffer)
	err := tmplt.Execute(buf, td)
	if err != nil {
		app.serverError(w, err)
		return
	}
	_, _ = buf.WriteTo(w)
}
