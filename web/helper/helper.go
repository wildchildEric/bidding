package helper

import (
	"fmt"
	"html/template"
	"wildchild.me/biddinginfo/db"
)

func Paginate(page db.Page) template.HTML {
	var totalPage int
	if page.TotalCount%page.CountPerPage == 0 {
		totalPage = page.TotalCount / page.CountPerPage
	} else {
		totalPage = page.TotalCount/page.CountPerPage + 1
	}
	pageInfo := fmt.Sprintf("<h5>共%d条 每页%d条  当前第%d页 共%d页</h5>", page.TotalCount, page.CountPerPage, page.CurrentPage, totalPage)

	if totalPage == 1 {
		str := `<div class="well">
				<div class="row">
					<div class="col-sm-12 col-md-4" style="text-align: left;font-size: 14px;">
						%s
					</div>
				</div>
			</div>`

		return template.HTML(fmt.Sprintf(str, pageInfo))
	}

	str := `<div class="well">
				<div class="row">
					<div class="col-sm-12 col-md-4" style="text-align: left;font-size: 14px;">
						%s
					</div>
					<div class="col-sm-12 col-md-8" style="text-align: right;">
						<ul class="pagination pagination-sm">%s</ul>
					</div>
				</div>
			</div>`

	pageHead := fmt.Sprintf(`<li><a href="?page=%d" aria-label="Previous"><span aria-hidden="true">« 首页</span></a></li>`, 1)
	pagePre := fmt.Sprintf(`<li><a href="?page=%d" aria-label="Previous"><span aria-hidden="true">‹ 上一页</span></a></li>`, page.CurrentPage-1)
	pageNext := fmt.Sprintf(`<li><a href="?page=%d" aria-label="Next"><span aria-hidden="true">下一页 ›</span></a></li>`, page.CurrentPage+1)
	pageEnd := fmt.Sprintf(`<li><a href="?page=%d" aria-label="Next"><span aria-hidden="true">尾页 »</span></a></li>`, totalPage)

	var pageLinks string
	linkNum := 5
	plinkActive := `<li class="active"><a>%d</a></li>`
	plink := `<li ><a href="?page=%d">%d</a></li>`
	plinkDisable := `<li class="disabled"><a href="#">…</a></li>`

	f := func(p db.Page, i int) string {
		if i == page.CurrentPage {
			return fmt.Sprintf(plinkActive, i)
		} else {
			return fmt.Sprintf(plink, i, i)
		}
	}
	if totalPage > linkNum {
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
	} else {
		for i := 1; i <= totalPage; i++ {
			pageLinks += f(page, i)
		}
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
}
