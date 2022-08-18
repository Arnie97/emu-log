package common_test

import (
	"fmt"

	"github.com/arnie97/emu-log/common"
)

func ExampleConf() {
	common.MockConf()

	fmt.Println(int64(common.Conf().Request.Interval))
	common.PrettyPrint(common.Conf().Request.Interval)

	serials := make(chan string)
	go func() {
		for _, a := range common.Conf().Adapters {
			for _, rule := range a.SearchSpace {
				rule.Emit(serials)
			}
		}
		close(serials)
	}()
	for s := range serials {
		fmt.Println(s)
	}

	// Output:
	// 2460
	// "2.46µs"
	// CRH5-001A
	// CRH3-002C
	// CRH3-004C
	// CRH2-001A
	// CRH2-002A
	// CRH2-003A
}
