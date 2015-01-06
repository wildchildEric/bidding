package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
