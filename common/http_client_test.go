package common_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/arnie97/emu-log/common"
)

func ExampleHTTPClient() {
	common.DisableMockHTTPClient()

	x := common.HTTPClient()
	y := common.HTTPClient()
	fmt.Println("x == y:  ", x == y)
	fmt.Println("x != nil:", x != nil)

	resp, err := x.Get("https://httpbin.org/user-agent")
	fmt.Println("get err: ", err)
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("read err:", err)

	s := struct {
		UserAgent string `json:"user-agent"`
	}{}
	err = json.Unmarshal(body, &s)
	fmt.Println("load err:", err)

	fmt.Println("ua equal:", s.UserAgent == common.UserAgent)

	// Output:
	// x == y:   true
	// x != nil: true
	// get err:  <nil>
	// read err: <nil>
	// load err: <nil>
	// ua equal: true
}
