package chinabidding

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
	"wildchild.me/biddinginfo/util"
)

const (
	ROOT_URL          = "http://www.chinabidding.com.cn"
	START_URL_DAILY   = ROOT_URL + "/search/searchzbw/search2?keywords=&areaid=7&categoryid=&b_date=day"
	START_URL_MONTHLY = ROOT_URL + "/search/searchzbw/search2?keywords=&areaid=7&categoryid=&b_date=month"
	START_URL_YEARLY  = ROOT_URL + "/search/searchzbw/search2?keywords=&areaid=7&categoryid=&b_date=year"
	LOGIN_PAGE_URL    = ROOT_URL + "/cblcn/member.login/login"
	LOGIN_CHECK_URL   = ROOT_URL + "/cblcn/member.login/logincheck"
	REQUEST_INTERVAL  = 10 * time.Millisecond
	REQUEST_TIME_OUT  = 3 * time.Second
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
	util.LogInvokeTime(func() {
		cookies, err := util.GetCookies(LOGIN_PAGE_URL)
		if err != nil {
			log.Fatal(err)
		}
		all_items := make([]*Item, 0, 4100)
		all_item_urls := make([]string, 0, 4100)
		html_str, err := util.GetPage(START_URL_MONTHLY, cookies)
		if err != nil {
			log.Println(err)
		}
		urls, err := ParseListPageToLinks(html_str)
		if err != nil {
			log.Println(err)
		}
		htmls := util.DownLoadPages(urls, cookies, REQUEST_INTERVAL, REQUEST_TIME_OUT)
		for i := 0; i < len(htmls); i++ {
			items, err := ParseListPageToItems(htmls[i])
			if err != nil {
				log.Println(err)
				continue
			}
			all_items = append(all_items, items...)
			item_urls := make([]string, 0, len(items))
			for _, v := range items {
				item_urls = append(item_urls, v.UrlDetail)
			}
			all_item_urls = append(all_item_urls, item_urls...)
		}

		// for i, item := range all_items {
		// 	fmt.Printf("%d %v\n", i, item)
		// }

		// for i, u := range all_item_urls {
		// 	fmt.Printf("%d %v\n", i, u)
		// }

		log.Println(len(all_items))
		log.Println(len(all_item_urls))

		// cookies, err = util.Login(LOGIN_CHECK_URL,
		// 	map[string]string{"name": "nmzb", "password": "NMzb2014"}, cookies)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// htmls = util.DownLoadPages(all_item_urls, cookies, REQUEST_INTERVAL, REQUEST_TIME_OUT)
		// log.Printf("%d", len(htmls))

		// for i, h := range htmls {
		// 	ParseDetailPage(all_items[i], h)
		// }

		// for i, item := range all_items {
		// 	fmt.Printf("%d %+v\n", i, item)
		// }

		// body, err := util.GetPage("http://www.chinabidding.com.cn/zbgg/F5hc.html", cookies)
		// item := &Item{}
		// ParseDetailPage(item, body)
		// fmt.Println(item.AgentName)
	})
}
