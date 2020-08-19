package common_test

import (
	"fmt"

	"github.com/arnie97/emu-log/common"
)

func ExampleCatch() {
	fmt.Println(panicFree())
	// Output: BOOM!
}

func panicFree() (err error) {
	defer common.Catch(&err)
	panic(fmt.Errorf("BOOM!"))
}
