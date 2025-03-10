package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", "Go")

	paths := []string{"./v23/ui/html/base.gohtml",
		"./v23/ui/html/partials/nav.gohtml",
		"./v23/ui/html/pages/home.gohtml",
	}

	tmplt, err := template.ParseFiles(paths...)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	err = tmplt.ExecuteTemplate(w, "base", nil)
	if err != nil {
		log.Print(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func snippetView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	_, err = fmt.Fprint(w, "Display a snippet with id: ", id)
	if err != nil {
		log.Fatal(err)
	}
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("Display a form for creating a new snippet"))
	if err != nil {
		log.Fatal(err)
	}
}

func snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
	_, err := w.Write([]byte("Creating a new snippet"))
	if err != nil {
		log.Fatal(err)
	}
}
