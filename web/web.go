package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"wildchild.me/biddinginfo/db"
	"wildchild.me/biddinginfo/web/helper"
)

const (
	templatesDir = "web/template/"
)

var (
	templatesMap map[string]*template.Template
	lk           sync.Mutex
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

	templateFuncs := template.FuncMap{
		"paginate": helper.Paginate,
	}

	for _, vl := range layouts {
		for _, vv := range views {
			lName := getPathName(vl)
			vName := getPathName(vv)
			key := fmt.Sprintf("%s_%s", lName, vName)
			t := template.New(key)
			t.Funcs(templateFuncs)
			t = template.Must(t.ParseFiles(vl, vv))
			templatesMap[key] = t
		}
	}
}

func render(w http.ResponseWriter, tmplName string, data interface{}) {
	lName := "application"
	key := fmt.Sprintf("%s_%s", lName, tmplName)
	lk.Lock()
	t, ok := templatesMap[key]
	lk.Unlock()
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
	pn, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		log.Println(err)
	}
	page, err := db.GetPage(pn, 5000)
	if err != nil {
		log.Println(err)
	}
	data := map[string]interface{}{"page": page}
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
