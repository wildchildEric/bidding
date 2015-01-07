package util

import (
	"errors"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	// "net/url"
	// "runtime"
	"strings"
	"time"
)

const (
	MAX_GORUTINE_NUM int = 1000
)

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
	_, err = client.Get(url)
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

func GetPage(urlStr string) (string, error) {
	// u, err := url.Parse(urlStr)
	// if err != nil {
	// 	return "", err
	// }
	// jar.SetCookies(u, cookies)
	client := http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// log.Printf("redirect %s to %s", urlStr, req.URL)
			// if len(via) >= 100 {
			return errors.New("stopped after 100 redirects ")
			// }
			// return nil
		},
	}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Response for %s with incorrect status code: %d", urlStr, resp.StatusCode))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func GetPageAsync(urlStr string) (<-chan string, <-chan string) {
	ch_content := make(chan string)
	ch_failed_url := make(chan string)
	go func() {
		html_str, err := GetPage(urlStr)
		if err != nil {
			log.Println(err)
			ch_failed_url <- urlStr
			return
		}
		ch_content <- html_str
	}()
	return ch_content, ch_failed_url
}

func DownLoadPages(urls []string, interval time.Duration, timeout time.Duration) []string {

	log.Printf("Downloading %d urls", len(urls))
	arr_chan := make([][2]<-chan string, 0, len(urls))
	htmls := make([]string, 0, len(urls))
	failed_urls := make([]string, 0, MAX_GORUTINE_NUM)

	for i, u := range urls {
		if i <= MAX_GORUTINE_NUM {
			time.Sleep(interval)
			ch0, ch1 := GetPageAsync(u)
			arr_chan = append(arr_chan, [2]<-chan string{ch0, ch1})
		} else {
			failed_urls = append(failed_urls, urls[i])
		}
	}
	for i, chan_arr := range arr_chan {
		ch0 := chan_arr[0]
		ch1 := chan_arr[1]
		select {
		case h := <-ch0:
			htmls = append(htmls, h)
		case u := <-ch1:
			failed_urls = append(failed_urls, u)
		case <-time.After(timeout):
			failed_urls = append(failed_urls, urls[i])
			log.Printf("%d item timed out: %s", i, urls[i])
		}
	}
	if len(failed_urls) > 0 {
		arr := DownLoadPages(failed_urls,
			interval+10*time.Millisecond,
			timeout+1*time.Second)
		htmls = append(htmls, arr...)
	}
	return htmls
}

func LogInvokeTime(f func()) {
	start := time.Now()
	f()
	log.Println(time.Now().Sub(start))
}
