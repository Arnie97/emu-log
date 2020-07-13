package common

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type mockHTTPClient struct {
	resp *http.Response
	err  error
}

func MockHTTPClientRespBody(body string) {
	confOnce.Do(func() {})
	httpOnce.Do(func() {})
	httpClient = &mockHTTPClient{resp: &http.Response{}}
	mockBody := ioutil.NopCloser(strings.NewReader(body))
	httpClient.(*mockHTTPClient).resp.Body = mockBody
}

func DisableMockHTTPClient() {
	confOnce = sync.Once{}
	httpOnce = sync.Once{}
	httpClient = &http.Client{
		Timeout:   RequestTimeout,
		Transport: &setDefaultHeaders{},
	}
}

func (x *mockHTTPClient) Do(*http.Request) (*http.Response, error) {
	return x.resp, x.err
}

func (x *mockHTTPClient) Get(url string) (*http.Response, error) {
	return x.resp, x.err
}

func (x *mockHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return x.resp, x.err
}

func (x *mockHTTPClient) PostForm(url string, data url.Values) (*http.Response, error) {
	return x.resp, x.err
}
