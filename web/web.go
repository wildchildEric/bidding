package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

var (
	templatesDir = "web/template/"
	templatesMap map[string]*template.Template
	templates    = template.Must(template.ParseGlob("web/template/*.html"))
)

func getPathName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(path)
	name := base[:strings.LastIndex(base, ext)]
	return name
}

func init() {
	if templatesMap == nil {
		templatesMap = make(map[string]*template.Template)
	}

	layouts, err := filepath.Glob(templatesDir + "layouts/*.html")
	if err != nil {
		log.Panic(err)
	}
	views, err := filepath.Glob(templatesDir + "views/*.html")
	if err != nil {
		log.Panic(err)
	}

	for _, vl := range layouts {
		for _, vv := range views {
			lName := getPathName(vl)
			vName := getPathName(vv)
			key := fmt.Sprintf("%s_%s", lName, vName)
			t := template.Must(template.New("").ParseFiles(vl, vv))
			templatesMap[key] = t
		}
	}
}

func render(w http.ResponseWriter, tmplName string) {
	lName := "base"
	key := fmt.Sprintf("%s_%s", lName, tmplName)
	t, ok := templatesMap[key]
	if !ok {
		http.Error(w, "No such template.", http.StatusInternalServerError)
	}
	err := t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	render(w, "root")
}

func Start() {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", rootHandler)
	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
