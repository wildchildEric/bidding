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
		"paginate": func(page db.Page) template.HTML {
			str := `<div class="row">
						<div class="col-sm-12 col-md-4" style="text-align: left;font-size: 14px;">%s</div>
						<div class="col-sm-12 col-md-8" style="text-align: right;">
					      <ul class="pagination pagination-sm">%s</ul>
						</div>
					</div>`
			var totalPage int
			if page.TotalCount%page.CountPerPage == 0 {
				totalPage = page.TotalCount / page.CountPerPage
			} else {
				totalPage = page.TotalCount/page.CountPerPage + 1
			}
			pageInfo := fmt.Sprintf("共%d条 每页%d条  当前第%d页 共%d页", page.TotalCount, page.CountPerPage, page.CurrentPage, totalPage)

			pageHead := fmt.Sprintf(`<li><a href="?page=%d" aria-label="Previous"><span aria-hidden="true">« 首页</span></a></li>`, 1)
			pagePre := fmt.Sprintf(`<li><a href="?page=%d" aria-label="Previous"><span aria-hidden="true">‹ 上一页</span></a></li>`, page.CurrentPage-1)
			pageNext := fmt.Sprintf(`<li><a href="?page=%d" aria-label="Next"><span aria-hidden="true">下一页 ›</span></a></li>`, page.CurrentPage+1)
			pageEnd := fmt.Sprintf(`<li><a href="?page=%d" aria-label="Next"><span aria-hidden="true">尾页 »</span></a></li>`, totalPage)

			var pageLinks string
			linkNum := 5
			plinkActive := `<li class="active"><a href="?page=%d">%d</a></li>`
			plink := `<li ><a href="?page=%d">%d</a></li>`
			plinkDisable := `<li class="disabled"><a href="#">…</a></li>`

			f := func(p db.Page, i int) string {
				if i == page.CurrentPage {
					return fmt.Sprintf(plinkActive, i, i)
				} else {
					return fmt.Sprintf(plink, i, i)
				}
			}

			if page.CurrentPage < linkNum {
				for i := 1; i <= linkNum; i++ {
					pageLinks += f(page, i)
				}
				pageLinks += plinkDisable
			} else if page.CurrentPage > totalPage-linkNum {
				pageLinks += plinkDisable
				for i := totalPage - linkNum; i <= totalPage; i++ {
					pageLinks += f(page, i)
				}
			} else {
				pageLinks += plinkDisable
				for i := page.CurrentPage - 2; i < page.CurrentPage+(linkNum-2); i++ {
					pageLinks += f(page, i)
				}
				pageLinks += plinkDisable
			}

			strLinks := `%s %s %s %s %s`
			if page.CurrentPage == 1 {
				strLinks = fmt.Sprintf(strLinks, "", "", pageLinks, pageNext, pageEnd)
			} else if page.CurrentPage == totalPage {
				strLinks = fmt.Sprintf(strLinks, pageHead, pagePre, pageLinks, "", "")
			} else {
				strLinks = fmt.Sprintf(strLinks, pageHead, pagePre, pageLinks, pageNext, pageEnd)
			}
			return template.HTML(fmt.Sprintf(str, pageInfo, strLinks))
		},
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
	page, err := db.GetPage(pn, 100)
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
