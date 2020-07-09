package common

import (
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	RequestTimeout = 10 * time.Second
	UserAgent      = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/7.0.8(0x17000820) NetType/4G Language/zh_CN"
)

var (
	httpClient httpRequester
	httpOnce   sync.Once
)

type (
	httpRequester interface {
		Do(*http.Request) (*http.Response, error)
		Get(url string) (*http.Response, error)
		Post(url, contentType string, body io.Reader) (*http.Response, error)
		PostForm(url string, data url.Values) (*http.Response, error)
	}
	setDefaultHeaders struct{}
)

func (setDefaultHeaders) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", UserAgent)
	return http.DefaultTransport.RoundTrip(req)
}

func HTTPClient() httpRequester {
	httpOnce.Do(func() {
		httpClient = &http.Client{
			Timeout:   RequestTimeout,
			Transport: &setDefaultHeaders{},
		}
	})
	return httpClient
}
