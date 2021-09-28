package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleHarbin_TrainNo() {
	printTrainNo(
		adapters.Harbin{},
		"harbin_basic.html",
	)
	// Output:
	//
	// false
	// "D28"          "2021-09-28"
}

func ExampleHarbin_VehicleNo() {
	printVehicleNo(adapters.Harbin{},
		"harbin_basic.html",
	)
	// Output:
	//
	// ""             false
}
