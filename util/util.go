package util

import (
	"errors"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type Page struct {
	Url     string
	Content string
}

var (
	jar *cookiejar.Jar
)

func InitCookieJar(url string) error {
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	var err error
	jar, err = cookiejar.New(&options)
	if err != nil {
		return err
	}
	client := &http.Client{
		Jar: jar,
	}
	_, err = client.Head(url)
	return err
}
func Login(url string, loginMap map[string]string) error {
	arr := make([]string, 0, 2)
	for k, v := range loginMap {
		arr = append(arr, fmt.Sprintf("%s=%s", k, v))
	}
	client := &http.Client{
		Jar: jar,
	}
	req, err := http.NewRequest(
		"POST",
		url,
		strings.NewReader(strings.Join(arr, "&")))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Login Failed.")
	}
	return nil
}

func GetPage(urlStr string) (*Page, error) {
	client := http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects ")
			}
			return nil
		},
	}

	resp, err := client.Get(urlStr)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Response for %s with incorrect status code: %d", urlStr, resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close() //instead of useing defer because: http://stackoverflow.com/questions/12952833/lookup-host-no-such-host-error-in-go
	if err != nil {
		return nil, err
	}
	return &Page{urlStr, string(body)}, nil
}

type ChanPair struct {
	ch0 <-chan *Page
	ch1 <-chan string
}

func GetPageAsync(urlStr string) ChanPair {
	ch_content := make(chan *Page)
	ch_failed_url := make(chan string)
	go func() {
		page, err := GetPage(urlStr)
		if err != nil {
			log.Println(err)
			ch_failed_url <- urlStr
			return
		}
		ch_content <- page
	}()
	return ChanPair{ch_content, ch_failed_url}
}

func DownLoadPages(urls []string, interval time.Duration, timeout time.Duration, beforeGetPage func(int)) []*Page {

	log.Printf("Downloading %d urls", len(urls))
	arr_chan := make([]ChanPair, 0, len(urls))
	pages := make([]*Page, 0, len(urls))
	failed_urls := make([]string, 0, len(urls)/2+1)

	for i, u := range urls {
		time.Sleep(interval)
		if beforeGetPage != nil {
			beforeGetPage(i)
		}
		chPair := GetPageAsync(u)
		arr_chan = append(arr_chan, chPair)
	}
	for i, chan_pair := range arr_chan {
		select {
		case p := <-chan_pair.ch0:
			pages = append(pages, p)
		case u := <-chan_pair.ch1:
			failed_urls = append(failed_urls, u)
		case <-time.After(timeout):
			failed_urls = append(failed_urls, urls[i])
			log.Printf("%d item timed out: %s", i, urls[i])
		}
	}
	if len(failed_urls) > 0 {
		arr := DownLoadPages(failed_urls,
			interval+10*time.Millisecond,
			timeout+1*time.Second,
			beforeGetPage)
		pages = append(pages, arr...)
	}
	return pages
}

func LogInvokeTime(f func()) {
	start := time.Now()
	f()
	log.Println(time.Now().Sub(start))
}
