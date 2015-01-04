package chinabidding

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	ROOT_URL          = "http://www.chinabidding.com.cn"
	START_URL_DAILY   = ROOT_URL + "/search/searchzbw/search2?keywords=&areaid=7&categoryid=&b_date=day"
	START_URL_MONTHLY = ROOT_URL + "/search/searchzbw/search2?keywords=&areaid=7&categoryid=&b_date=month"
	START_URL_YEARLY  = ROOT_URL + "/search/searchzbw/search2?keywords=&areaid=7&categoryid=&b_date=year"
	LOGIN_PAGE_URL    = ROOT_URL + "/cblcn/member.login/login"
	LOGIN_CHECK_URL   = ROOT_URL + "/cblcn/member.login/logincheck"
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

func GetCookies(url string) ([]*http.Cookie, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, err
	}
	return resp.Cookies(), nil
}

func Login(name string, pass string, cookies []*http.Cookie) ([]*http.Cookie, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		LOGIN_CHECK_URL,
		strings.NewReader(fmt.Sprintf("name=%s&password=%s", name, pass)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for _, c := range cookies {
		req.AddCookie(c)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 200 {
		return resp.Cookies(), nil
	} else {
		return nil, errors.New("Login Failed.")
	}
}

func GetPage(urlStr string, cookies []*http.Cookie) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func ParseListPageToItems(html_string string) ([]*Item, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html_string))
	if err != nil {
		return nil, err
	}
	var items []*Item
	parse_func := func(i int, s *goquery.Selection) {
		item := &Item{}
		s.Find("td").Each(func(i int, s *goquery.Selection) {
			switch i {
			case 1:
				item.Title = strings.TrimSpace(s.Text())
				href, exist := s.Find("a").Attr("href")
				if exist {
					item.UrlDetail = ROOT_URL + href
				}
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
		items = append(items, item)
	}
	doc.Find(".listrow1").Each(parse_func)
	doc.Find(".listrow2").Each(parse_func)
	return items, nil
}

func ParseDetailPage(item *Item, html_string string) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html_string))
	if err != nil {
		return err
	}
	doc.Find(".f_l.nr_bt1_sf.f_12 li").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "招标代理") {
			splited := strings.Split(s.Text(), ":")
			agent := splited[len(splited)-1]
			item.AgentName = agent
		}
	})
	return nil
}

func ParseListPageToLinks(html_string string) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html_string))
	if err != nil {
		return nil, err
	}
	s := doc.Find("#pages a").Last()
	href, exist := s.Attr("href")
	if !exist {
		return nil, errors.New("not end page link found")
	}
	u, err := url.Parse(ROOT_URL + href)
	if err != nil {
		return nil, err
	}
	m, _ := url.ParseQuery(u.RawQuery)
	max_page, err := strconv.Atoi(m["page"][0])
	if err != nil {
		return nil, err
	}
	list_page_urls := make([]string, 0, max_page)
	for i := 1; i <= max_page; i++ {
		m.Set("page", strconv.Itoa(i))
		url := u.Scheme + "://" + u.Host + u.Path + "?" + m.Encode()
		list_page_urls = append(list_page_urls, url)
	}
	return list_page_urls, nil
}

func Start() {
	cookies, err := GetCookies(LOGIN_PAGE_URL)
	if err != nil {
		log.Fatal(err)
	}
	// cookies, err = Login("nmzb", "NMzb2014", cookies)
	// body := GetPage("http://www.chinabidding.com.cn/zbgg/F5hc.html", cookies)
	// item := &Item{}
	// ParseDetailPage(item, body)
	// fmt.Println(item.AgentName)

	all_items := make([]*Item, 0, 4100)
	list_html_str, err := GetPage(START_URL_MONTHLY, cookies)
	if err != nil {
		log.Fatal(err)
	}
	url_list, err := ParseListPageToLinks(list_html_str)
	if err != nil {
		log.Fatal(err)
	}
	for i, u := range url_list {
		html_str, err := GetPage(u, cookies)
		if err != nil {
			log.Fatal(err)
		}
		items, err := ParseListPageToItems(html_str)
		if err != nil {
			log.Fatal(err)
		}
		all_items = append(all_items, items...)
		fmt.Printf("%d %d all_items length: %d cap: %d\n", i, len(items), len(all_items), cap(all_items))
	}

	for i, item := range all_items {
		fmt.Printf("%d, %v\n", i, item)
	}
}
