package common

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	RequestTimeout = 20 * time.Second
)

type (
	HTTPRequester interface {
		Do(*http.Request) (*http.Response, error)
		Get(url string) (*http.Response, error)
		Post(url, contentType string, body io.Reader) (*http.Response, error)
		PostForm(url string, data url.Values) (*http.Response, error)
	}
)

func (conf *RequestConf) RoundTrip(req *http.Request) (*http.Response, error) {
	if conf.Interval > 0 {
		time.Sleep(time.Duration(conf.Interval))
	}
	if len(conf.UserAgent) > 0 {
		req.Header.Set("user-agent", conf.UserAgent)
	}
	return http.DefaultTransport.RoundTrip(req)
}

func SetCookies(req *http.Request, cookies []*http.Cookie) {
	if len(cookies) == 0 {
		return
	}

	var buf bytes.Buffer
	for _, each := range cookies {
		buf.WriteString(each.Name)
		buf.WriteRune('=')
		buf.WriteString(each.Value)
		buf.WriteString("; ")
	}
	req.Header.Set("cookie", buf.String()[:buf.Len()-2])
}

func HTTPClient(roundTripper ...http.RoundTripper) HTTPRequester {
	if mockHTTPClientInstance != nil {
		return mockHTTPClientInstance
	}
	if roundTripper == nil {
		roundTripper = []http.RoundTripper{Conf().Request}
	}
	return &http.Client{
		Timeout:   RequestTimeout,
		Transport: roundTripper[0],
	}
}
