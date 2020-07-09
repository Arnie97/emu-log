package common

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type mockHTTPClient struct {
	resp *http.Response
	err  error
}

func MockHTTPClient() httpRequester {
	confOnce.Do(func() {})
	httpOnce.Do(func() {
		httpClient = &mockHTTPClient{resp: &http.Response{}}
	})
	return httpClient
}

func SetMockHTTPClientRespBody(body string) {
	mockBody := ioutil.NopCloser(strings.NewReader(body))
	MockHTTPClient().(*mockHTTPClient).resp.Body = mockBody
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
