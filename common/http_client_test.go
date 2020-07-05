package common_test

import (
	"testing"

	"github.com/arnie97/emu-log/common"
	"github.com/stretchr/testify/assert"
)

func TestHTTPClient(t *testing.T) {
	x := common.HTTPClient()
	y := common.HTTPClient()
	assert.Equal(t, x, y)
}
