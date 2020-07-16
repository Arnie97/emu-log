package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleGuangzhou_TrainNo() {
	printTrainNo(
		adapters.Guangzhou{},
		"guangzhou_full.json",
		"guangzhou_invalid.json",
	)
	// Output:
	// "G1363"        false ""
	// ""             true  ""
}

func ExampleGuangzhou_VehicleNo() {
	printVehicleNo(
		adapters.Guangzhou{},
		"guangzhou_full.json",
		"guangzhou_invalid.json",
	)
	// Output:
	// "CR8+8-11876@" false
	// ""             true
}
