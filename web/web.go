package web

import (
	// "fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	templates = template.Must(template.ParseFiles("web/template/application.html"))
)

func renderTemplate(w http.ResponseWriter, tmpl string) {
	err := templates.ExecuteTemplate(w, tmpl+".html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "application")
}

func Start() {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", rootHandler)
	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
