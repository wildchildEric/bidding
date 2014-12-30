package main

import (
	"fmt"
	"wildchild.me/biddinginfo/spiders/chinabidding"
)

func main() {
	cookies := chinabidding.GetCookies("http://www.chinabidding.com.cn/cblcn/member.login/login")
	// cookies = chinabidding.Login("nmzb", "NMzb2014", cookies)
	// body := getPage("http://www.chinabidding.com.cn/zbgs/jMGQ.html", cookies)
	// fmt.Println(body)

	list_html_str := chinabidding.GetPage("http://www.chinabidding.com.cn/search/searchzbw/search2?keywords=&areaid=7&categoryid=&b_date=month", cookies)
	items := chinabidding.ParseListPage(list_html_str)
	for _, ele := range items {
		fmt.Println(ele)
		// fmt.Printf("%d %q %q %q %q %q %q\n", index, ele.Title, ele.Category, ele.Region, ele.Industry, ele.Date, ele.UrlDetail)
	}
}
