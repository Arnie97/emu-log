package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleWuhan_TrainNo() {
	printTrainNo(
		adapters.Wuhan{},
		"wuhan_full.html",
		"wuhan_basic.html",
		"wuhan_invalid.html",
	)
	// Output:
	// "G1730"        false ""
	// ""             true  ""
	// ""             true  ""
}

func ExampleWuhan_VehicleNo() {
	printVehicleNo(
		adapters.Wuhan{},
		"wuhan_full.html",
		"wuhan_basic.html",
		"wuhan_invalid.html",
	)
	// Output:
	// "CR400AF2151"  false
	// "CRH380D1545"  false
	// ""             true
}
