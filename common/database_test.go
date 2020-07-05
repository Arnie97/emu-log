package common_test

import (
	"testing"

	"github.com/arnie97/emu-log/common"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	x := common.DB()
	y := common.DB()
	assert.Equal(t, x, y)
}
