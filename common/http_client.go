package common

import (
	"net/http"
	"sync"
	"time"
)

const (
	RequestTimeout = 10 * time.Second
	UserAgent      = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/7.0.8(0x17000820) NetType/4G Language/zh_CN"
)

var (
	httpClient *http.Client
	httpOnce   sync.Once
)

type setDefaultHeaders struct{}

func (setDefaultHeaders) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", UserAgent)
	return http.DefaultTransport.RoundTrip(req)
}

func HTTPClient() *http.Client {
	httpOnce.Do(func() {
		httpClient = &http.Client{
			Timeout:   RequestTimeout,
			Transport: &setDefaultHeaders{},
		}
	})
	return httpClient
}
