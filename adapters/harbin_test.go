package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleHarbin_TrainNo() {
	printTrainNo(
		adapters.Harbin{},
		"harbin_basic.html",
		"harbin_invalid.html",
	)
	// Output:
	//
	// false
	// "G1206/7"      ""
	//
	// true
}

func ExampleHarbin_VehicleNo() {
	printVehicleNo(adapters.Harbin{},
		"harbin_basic.html",
		"harbin_invalid.html",
	)
	// Output:
	//
	// "CRH380BG@42s" false
	// ""             true
}
