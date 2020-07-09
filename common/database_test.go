package common_test

import (
	"fmt"

	"github.com/arnie97/emu-log/common"
)

func ExampleDB() {
	x := common.DB()
	y := common.DB()
	_, err := x.Exec(`SELECT 1;`)
	common.CountRecords("emu_log")

	fmt.Println("x == y:  ", x == y)
	fmt.Println("x != nil:", x != nil)
	fmt.Println("error:   ", err)
	// Output:
	// x == y:   true
	// x != nil: true
	// error:    <nil>
}
