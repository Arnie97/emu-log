package common_test

import (
	"fmt"
	"os"

	"github.com/arnie97/emu-log/common"
)

func mockConf() {
	file, err := os.Create(common.AppPath() + common.ConfPath)
	if err != nil {
		fmt.Println(err)
	}
	file.Write([]byte(`
		[request]
		interval = "0.00246ms"
		user-agent = "Mozilla/5.0"

		[[adapters.X.search]]
		format = "CRH5-%03dA"

		[[adapters.X.search]]
		format = "CRH3-%03dC"
		min = 2
		max = 4
		step = 2

		[[adapters.X.search]]
		format = "CRH2-%03dA"
		max = 3
	`))
	file.Close()
}

func ExampleConf() {
	mockConf()

	fmt.Println(int64(common.Conf().Request.Interval))
	common.PrettyPrint(common.Conf().Request.Interval)

	serials := make(chan string)
	go func() {
		for _, rule := range common.Conf().Adapters["X"].SearchSpace {
			rule.Emit(serials)
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
