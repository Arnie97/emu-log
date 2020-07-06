package common_test

import (
	"testing"

	"github.com/arnie97/emu-log/common"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	x := common.DB()
	y := common.DB()
	assert.NotNil(t, x)
	assert.Equal(t, x, y)
	common.DB().Exec(`SELECT 1;`)
	common.CountRecords("emu_log")
}
