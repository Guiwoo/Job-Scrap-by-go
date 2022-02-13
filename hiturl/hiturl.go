package hiturl

import (
	"errors"
	"fmt"
	"net/http"
)

type result struct {
	url    string
	status string
}

var errRequestFailed = errors.New("Request Failed")

func hitUrls(url string, c chan<- result) {
	fmt.Println("Checking: ", url)
	res, err := http.Get(url)
	status := "Okay"
	if err != nil || res.StatusCode >= 400 {
		status = "FAILED"
	}
	c <- result{url: url, status: status}

}

func main() {
	results := make(map[string]string)
	c := make(chan result)
	urls := []string{
		"https://www.google.com/",
		"https://www.naver.com/",
		"https://www.facebook.com/",
		"https://www.amazon.com/",
		"https://www.daum.net/",
		"https://www.youtube.com/",
	}

	for _, url := range urls {
		go hitUrls(url, c)
	}

	for i := 0; i < len(urls); i++ {
		req_result := <-c
		results[req_result.url] = req_result.status
	}
	for url, status := range results {
		fmt.Println(url, status)
	}
}
