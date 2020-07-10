package common_test

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/arnie97/emu-log/common"
)

func ExampleMockHTTPClientRespBody() {
	common.MockHTTPClientRespBody("CRH6A-4002")
	x := common.HTTPClient()

	a, e := x.Do(nil)
	b, f := x.Get("")
	c, g := x.Post("", "", bytes.NewReader(nil))
	d, h := x.PostForm("", nil)
	fmt.Println("same resp:", a == b && b == c && c == d)

	body, i := ioutil.ReadAll(d.Body)
	j := d.Body.Close()
	fmt.Println("no errors:", e == f && f == g && g == h && h == i && i == j)
	fmt.Println("resp body:", string(body))

	// Output:
	// same resp: true
	// no errors: true
	// resp body: CRH6A-4002
}
