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
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Response with incorrect status code: %d", resp.StatusCode))
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

func DownLoadPages(urlList []string, cookies []*http.Cookie, requestInterval time.Duration, requestTimeout time.Duration) []string {
	ch := make(chan string)
	ch_f := make(chan string)
	arr_html := make([]string, 0, 4100)
	arr_failed_url := make([]string, 0, 100)
	for _, u := range urlList {
		time.Sleep(requestInterval)
		go func() {
			html_str, err := GetPage(u, cookies)
			if err != nil {
				log.Println(err)
				ch_f <- u
				return
			}
			ch <- html_str
		}()
	}
	for i := 0; i < len(urlList); i++ {
		timeout := time.After(requestTimeout)
		select {
		case html_str := <-ch:
			arr_html = append(arr_html, html_str)
		case failed_url := <-ch_f:
			arr_failed_url = append(arr_failed_url, failed_url)
		case <-timeout:
			log.Printf("%d item timed out", i)
			arr_failed_url = append(arr_failed_url, urlList[i])
		}
	}
	if len(arr_failed_url) > 0 {
		arr := DownLoadPages(arr_failed_url,
			cookies,
			requestInterval+10*time.Millisecond,
			requestTimeout+1*time.Second)
		arr_html = append(arr_html, arr...)
	}
	return arr_html
}

func LogInvokeTime(f func()) {
	start := time.Now()
	f()
	log.Println(time.Now().Sub(start))
}
