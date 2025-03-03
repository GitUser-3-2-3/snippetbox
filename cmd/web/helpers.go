package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
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

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	if td == nil {
		td = &templateData{}
	}
	td.CurrentYear = time.Now().Year()
	return td
}

//goland:noinspection GoUnusedParameter
func (app *application) render(w http.ResponseWriter, r *http.Request, name string, td *templateData) {
	tmplt, ok := app.templateCache[name]
	if !ok {
		err := fmt.Errorf("used template %s does not exist", name)
		app.serverError(w, err)
		return
	}
	buf := new(bytes.Buffer)
	err := tmplt.Execute(buf, app.addDefaultData(td, r))
	if err != nil {
		app.serverError(w, err)
		return
	}
	_, _ = buf.WriteTo(w)
}
