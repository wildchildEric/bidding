package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"wildchild.me/biddinginfo/db"
)

var (
	templatesDir = "web/template/"
	templatesMap map[string]*template.Template
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
			t := template.Must(template.New(key).ParseFiles(vl, vv))
			templatesMap[key] = t
		}
	}
}

func render(w http.ResponseWriter, tmplName string, data interface{}) {
	lName := "application"
	key := fmt.Sprintf("%s_%s", lName, tmplName)
	t, ok := templatesMap[key]
	if !ok {
		http.Error(w, "No such template.", http.StatusInternalServerError)
		return
	}
	err := t.ExecuteTemplate(w, lName, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handlerICon(w http.ResponseWriter, r *http.Request) {}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		log.Println(err)
	}
	items, err := db.GetAll(page)
	if err != nil {
		log.Println(err)
	}
	data := map[string]interface{}{"items": items, "count": len(items)}

	render(w, "root", data)
}

func Start() {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/favicon.ico", handlerICon)
	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
