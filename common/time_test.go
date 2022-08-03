package common_test

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/arnie97/emu-log/common"
)

func ExampleUnixMilli() {
	// The doomsday of gophers!
	t, err := time.Parse(time.RFC1123, "Sat, 12 Apr 2262 07:47:16 CST")
	fmt.Println(common.UnixMilli(t), err)

	fmt.Println(common.UnixMilli(time.Now()) == common.UnixMilli())

	defer func() {
		fmt.Println(recover() != nil)
	}()
	common.UnixMilli(time.Now(), time.Now())

	// Output:
	// 9223372036000 <nil>
	// true
	// true
}

func ExampleMockStaticUnixMilli() {
	mockTime := int64(rand.Uint32())
	common.MockStaticUnixMilli(mockTime)
	fmt.Println(common.UnixMilli() == mockTime)

	// Output: true
}
