package web

import (
	"bytes"
	"encoding/csv"
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

func convertCsvArray(items []*db.Item) [][]string {
	records := make([][]string, 0, len(items))
	records = append(records, []string{"标题", "类别", "地区", "行业", "日期", "招标代理", "链接"})
	for _, it := range items {
		r := make([]string, 0, 7)
		r = append(r, it.Title)
		r = append(r, it.Category)
		r = append(r, it.Region)
		r = append(r, it.Industry)
		r = append(r, it.Date)
		r = append(r, it.AgentName)
		r = append(r, it.UrlDetail)
		records = append(records, r)
	}
	return records
}

func handlerICon(w http.ResponseWriter, r *http.Request) {}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	pn, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		pn = 0
	}
	page, err := db.GetPage(pn, 200)
	if err != nil {
		log.Println(err)
	}
	data := map[string]interface{}{"page": page}
	render(w, "root", data)
}

func bulkActionHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	action := r.FormValue("bulk_action")
	ids := r.Form["model_ids"]

	switch action {
	case "bulk_export_select_selected":
		items, err := db.GetItems(ids)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		records := convertCsvArray(items)
		b := &bytes.Buffer{}
		wr := csv.NewWriter(b)
		wr.WriteAll(records)
		wr.Flush()
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment;filename=招标数据.csv")
		w.Write(b.Bytes())
	case "bulk_export_excluding_selected":
		fmt.Fprintf(w, "not implemented yet.")
	}
}

func Start() {
	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/bulk_action", bulkActionHandler)
	http.HandleFunc("/favicon.ico", handlerICon)
	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
