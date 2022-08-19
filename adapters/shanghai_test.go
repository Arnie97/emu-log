package adapters_test

import (
	"fmt"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
)

const (
	shanghaiTestSerial = "PQC47BDAA1B1AC46C2AD545D93B9E30BA3"
)

func ExampleShanghai_SerialEncrypt() {
	common.MockStaticUnixMilli(1632636012773)
	fmt.Println(string(adapters.Shanghai{}.SerialEncrypt(shanghaiTestSerial)))
	// Output: {"data":"bwmZtecxmGBJrIVOpsB8/n66ix922um4AhzjOb5eFuZKWQdSJzA1C0BJIFo5iv9C4QyafxmIswLZWx6AhA0szfxxP52mFA7xifBIS/66xhEuEeNgJTisY69iXu9WKtHYVHZ5ywMpI7sBr1Xu7BzmL50UmcLitNYlcWg8OC7ry+wTJqlRwbHeYk8zwia54qiQIBjZeKJITcPtM7c8midgnA=="}
}

func ExampleShanghai_Signature() {
	fmt.Println(adapters.Shanghai{}.Signature(map[string]string{
		"pqCode": shanghaiTestSerial,
	}))
	// Output: 2f9affaf878b65cf3a80
}

func ExampleShanghai_TrainNo() {
	printTrainNo(
		adapters.Shanghai{},
		"shanghai_full.json",
		"shanghai_basic.json",
		"shanghai_invalid.json",
	)
	// Output:
	//
	// false
	// "G8"           "2021-09-27 13:27:00"
	//
	// true
	//
	// true
}

func ExampleShanghai_UnitNo() {
	printUnitNo(adapters.Shanghai{},
		"shanghai_full.json",
		"shanghai_basic.json",
		"shanghai_invalid.json",
	)
	// Output:
	// H "CR400BFB5097" false false false
	// H ""              true  true false
	// H ""              true  true false
}
