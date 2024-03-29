package common

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func MockConf() {
	os.Link(
		filepath.Join("../adapters/testdata", confFile),
		confPath(),
	)
	Conf()
}

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
	mockHTTPClientInstance = &mockHTTPClient{body: body}
}

func MockHTTPClientRespBodyFromFile(mockFile string) {
	MockHTTPClientRespBody(string(ReadMockFile(mockFile)))
}

func MockHTTPClientError(err error) {
	mockHTTPClientInstance = &mockHTTPClient{err: err}
}

func DisableMockHTTPClient() {
	mockHTTPClientInstance = nil
}

func (x *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	resp := &http.Response{
		Request: req,
		Header:  http.Header{"Set-Cookie": {"JSESSIONID=1234"}},
	}
	if x.err == nil {
		resp.Body = ioutil.NopCloser(strings.NewReader(x.body))
	}
	return resp, x.err
}

func (x *mockHTTPClient) Get(url string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return x.Do(req)
}

func (x *mockHTTPClient) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return x.Do(req)
}

func (x *mockHTTPClient) PostForm(url string, data url.Values) (*http.Response, error) {
	return x.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
