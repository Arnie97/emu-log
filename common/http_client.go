package common

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	RequestInterval = 3 * time.Second
	RequestTimeout  = 10 * time.Second
	UserAgentWeChat = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 MicroMessenger/7.0.8(0x17000820) NetType/4G Language/zh_CN"
	UserAgentJDPay  = "Mozilla/5.0 (Linux; Android 7.1.2; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/67.0.3396.87 Mobile Safari/537.36/application=JDJR-App&clientType=android&#@jdPaySDK*#@"
)

type (
	httpRequester interface {
		Do(*http.Request) (*http.Response, error)
		Get(url string) (*http.Response, error)
		Post(url, contentType string, body io.Reader) (*http.Response, error)
		PostForm(url string, data url.Values) (*http.Response, error)
	}
	IntervalTransport struct{}
)

func (IntervalTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	time.Sleep(RequestInterval)
	req.Header.Set("user-agent", UserAgentWeChat)
	return http.DefaultTransport.RoundTrip(req)
}

func HTTPClient(roundTripper ...http.RoundTripper) httpRequester {
	if mockHTTPClientInstance != nil {
		return mockHTTPClientInstance
	}
	if roundTripper == nil {
		roundTripper = []http.RoundTripper{IntervalTransport{}}
	}
	return &http.Client{
		Timeout:   RequestTimeout,
		Transport: roundTripper[0],
	}
}
