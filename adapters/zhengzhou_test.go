package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleZhengzhou_TrainNo() {
	printTrainNo(
		&adapters.Zhengzhou{},
		"zhengzhou_full.json",
		"zhengzhou_basic.json",
		"zhengzhou_invalid.json",
	)
	// Output:
	//
	// false
	// "D2751"        "2020-07-12"
	//
	// true
	//
	// true
}

func ExampleZhengzhou_VehicleNo() {
	printVehicleNo(
		&adapters.Zhengzhou{},
		"zhengzhou_full.json",
		"zhengzhou_basic.json",
		"zhengzhou_invalid.json",
	)
	// Output:
	// "CRH5G5194"    false
	// "CRH380B3667"  false
	// ""             true
}
