package common

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"
)

func ReadMockFile(mockFile string) (content []byte) {
	content, err := ioutil.ReadFile(filepath.Join("testdata", mockFile))
	Must(err)
	return content
}

type mockHTTPClient struct {
	body string
	err  error
}

var mockHTTPClientInstance *mockHTTPClient

func MockHTTPClientRespBody(body string) {
	confOnce.Do(func() {})
	mockHTTPClientInstance = &mockHTTPClient{body, nil}
}

func MockHTTPClientRespBodyFromFile(mockFile string) {
	MockHTTPClientRespBody(string(ReadMockFile(mockFile)))
}

func DisableMockHTTPClient() {
	confOnce = sync.Once{}
	mockHTTPClientInstance = nil
}

func (x *mockHTTPClient) Do(*http.Request) (*http.Response, error) {
	mockBody := ioutil.NopCloser(strings.NewReader(x.body))
	resp := &http.Response{Body: mockBody}
	return resp, x.err
}

func (x *mockHTTPClient) Get(url string) (*http.Response, error) {
	return x.Do(nil)
}

func (x *mockHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	return x.Do(nil)
}

func (x *mockHTTPClient) PostForm(url string, data url.Values) (*http.Response, error) {
	return x.Do(nil)
}
