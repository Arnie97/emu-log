package common_test

import (
	"fmt"
	"os"

	"github.com/arnie97/emu-log/common"
)

func ExampleConf() {
	file, err := os.Create(common.AppPath() + "/emu-log.json")
	if err != nil {
		fmt.Println(err)
	}
	file.Write([]byte(`{"a": "1234", "b": "5678"}`))
	file.Close()

	fmt.Println(common.Conf("a"), common.Conf("b"))
	// Output:
	// 1234 5678
}
