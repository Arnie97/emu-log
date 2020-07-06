package common_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/arnie97/emu-log/common"
	"github.com/stretchr/testify/assert"
)

func TestHTTPClient(t *testing.T) {
	x := common.HTTPClient()
	y := common.HTTPClient()
	assert.NotNil(t, x)
	assert.Equal(t, x, y)

	resp, err := x.Get("https://httpbin.org/user-agent")
	assert.Nil(t, err)
	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	s := struct {
		UserAgent string `json:"user-agent"`
	}{}
	err = json.Unmarshal(body, &s)
	assert.Nil(t, err)
	assert.Equal(t, s.UserAgent, common.UserAgent)
}
