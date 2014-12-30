package main

import (
	"fmt"
	"wildchild.me/biddinginfo/spiders/chinabidding"
)

func main() {
	cookies := chinabidding.GetCookies(chinabidding.LOGIN_PAGE_URL)
	// cookies = chinabidding.Login("nmzb", "NMzb2014", cookies)
	// body := getPage("http://www.chinabidding.com.cn/zbgs/jMGQ.html", cookies)
	// fmt.Println(body)

	list_html_str := chinabidding.GetPage(chinabidding.START_URL_MONTHLY, cookies)
	items := chinabidding.ParseListPage(list_html_str)
	for i, ele := range items {
		fmt.Printf("%d %q %q %q %q %q %q\n", i, ele.Title, ele.Category, ele.Region, ele.Industry, ele.Date, ele.UrlDetail)
	}
}
