package common_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/arnie97/emu-log/common"
)

func ExampleHTTPClient() {
	common.DisableMockHTTPClient()

	x := common.HTTPClient()
	y := common.HTTPClient(http.DefaultTransport)
	fmt.Println("x == y:  ", x == y)

	const api = "https://httpbin.org/anything"
	req, _ := http.NewRequest(http.MethodPut, api, nil)
	common.SetCookies(req, nil)
	common.SetCookies(req, []*http.Cookie{
		{Name: "model", Value: "CRH6A"},
		{Name: "serial", Value: "4002"},
	})
	cookies := req.Header.Get("cookie")
	fmt.Println("cookies: ", cookies)

	resp, err := x.Do(req)
	fmt.Println("get err: ", err)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("read err:", err)

	s := struct {
		Headers struct {
			Cookies   string `json:"cookie"`
			UserAgent string `json:"user-agent"`
		} `json:"headers"`
	}{}
	err = json.Unmarshal(body, &s)
	fmt.Println("load err:", err)

	fmt.Println("ua equal:", s.Headers.UserAgent == common.UserAgentWeChat)
	fmt.Println("cookies: ", s.Headers.Cookies == cookies)

	// Output:
	// x == y:   false
	// cookies:  model=CRH6A; serial=4002
	// get err:  <nil>
	// read err: <nil>
	// load err: <nil>
	// ua equal: true
	// cookies:  true
}
