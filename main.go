package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Item struct {
	Title     string
	Category  string
	Region    string
	Industry  string
	Date      string
	AgentName string
	UrlDetail string
}

func getCookies(url string) []*http.Cookie {
	resp, err := http.Head(url)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Cookies()
}

func login(name string, pass string, cookies []*http.Cookie) []*http.Cookie {
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		"http://www.chinabidding.com.cn/cblcn/member.login/logincheck",
		strings.NewReader(fmt.Sprintf("name=%s&password=%s", name, pass)))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 200 {
		return resp.Cookies()
	} else {
		log.Fatal("Login failed.")
		return nil
	}
}

func getPage(urlStr string, cookies []*http.Cookie) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(body)
}

func parseListPage(html_string string) {
	reader := strings.NewReader(html_string)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		log.Fatal(err)
	}
	parse_func := func(i int, s *goquery.Selection) {
		item := &Item{}
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 1:
				item.Title = strings.TrimSpace(s.Text())
				item.UrlDetail, _ = s.Find("a").Attr("href")
			case 2:
				item.Category = strings.TrimSpace(s.Text())
			case 3:
				item.Region = strings.TrimSpace(s.Text())
			case 4:
				item.Industry = strings.TrimSpace(s.Text())
			case 5:
				slice_strings := strings.Split(strings.TrimSpace(s.Text()), "\n")
				item.Date = strings.TrimSpace(slice_strings[len(slice_strings)-1])
			}
		})
		fmt.Printf("%d %q %q %q %q %q %q\n", i, item.Title, item.Category, item.Region, item.Industry, item.Date, item.UrlDetail)
	}

	doc.Find(".listrow1").Each(parse_func)
	doc.Find(".listrow2").Each(parse_func)
}

func main() {
	cookies := getCookies("http://www.chinabidding.com.cn/cblcn/member.login/login")
	// cookies = login("nmzb", "NMzb2014", cookies)

	// body := getPage("http://www.chinabidding.com.cn/zbgs/jMGQ.html", cookies)
	// fmt.Println(body)

	list_html_str := getPage("http://www.chinabidding.com.cn/search/searchzbw/search2?keywords=&areaid=7&categoryid=&b_date=month", cookies)
	parseListPage(list_html_str)
}
