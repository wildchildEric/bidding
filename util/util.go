package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetCookies(url string) ([]*http.Cookie, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, err
	}
	return resp.Cookies(), nil
}

func GetPage(url string, cookies []*http.Cookie) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
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
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Response for %s with incorrect status code: %d", url, resp.StatusCode))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func Login(url string, loginMap map[string]string, cookies []*http.Cookie) ([]*http.Cookie, error) {
	arr := make([]string, 0, 2)
	for k, v := range loginMap {
		arr = append(arr, fmt.Sprintf("%s=%s", k, v))
	}
	client := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		url,
		strings.NewReader(strings.Join(arr, "&")))
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
	if resp.StatusCode == http.StatusOK {
		return resp.Cookies(), nil
	} else {
		return nil, errors.New("Login Failed.")
	}
}

func GetPageAsync(urlStr string, cookies []*http.Cookie) (<-chan string, <-chan string) {
	ch_content := make(chan string)
	ch_failed_url := make(chan string)
	go func() {
		html_str, err := GetPage(urlStr, cookies)
		if err != nil {
			log.Println(err)
			ch_failed_url <- urlStr
			return
		}
		ch_content <- html_str
	}()
	return ch_content, ch_failed_url
}

func DownLoadPages(urls []string, cookies []*http.Cookie,
	interval time.Duration, timeout time.Duration) []string {

	arr_chan := make([][2]<-chan string, 0, len(urls))
	htmls := make([]string, 0, len(urls))
	failed_urls := make([]string, 0, len(urls)/2+1)

	for _, u := range urls {
		time.Sleep(interval)
		ch0, ch1 := GetPageAsync(u, cookies)
		arr_chan = append(arr_chan, [2]<-chan string{ch0, ch1})
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
			log.Printf("%d item timed out %d", i, urls[i])
		}
	}
	if len(failed_urls) > 0 {
		arr := DownLoadPages(failed_urls,
			cookies,
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
