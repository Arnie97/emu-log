package common_test

import (
	"bufio"
	"bytes"
	"fmt"

	"github.com/arnie97/emu-log/common"
)

func ExampleMockHTTPClientRespBodyFromFile() {
	common.MockHTTPClientRespBodyFromFile("../http_client_mock_test.go")
	x := common.HTTPClient()

	a, e := x.Do(nil)
	b, f := x.Get("")
	c, g := x.Post("", "", bytes.NewReader(nil))
	d, h := x.PostForm("", nil)
	fmt.Println("same resp:", a == b || b == c || c == d)

	body, i := bufio.NewReader(d.Body).ReadString('\n')
	j := d.Body.Close()
	fmt.Println("no errors:", isAllNil(e, f, g, h, i, j))
	fmt.Print("resp body: ", body)

	common.MockHTTPClientError(fmt.Errorf("my sample error"))
	k, m := common.HTTPClient().Do(nil)
	fmt.Printf("err  mock: %v, %v", k.Body, m)

	// Output:
	// same resp: false
	// no errors: true
	// resp body: package common_test
	// err  mock: <nil>, my sample error
}

func isAllNil(values ...interface{}) bool {
	for _, v := range values {
		if v != nil {
			return false
		}
	}
	return true
}
