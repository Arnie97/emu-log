package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleKunming_TrainNo() {
	printTrainNo(
		&adapters.Kunming{},
		"kunming_full.json",
		"kunming_invalid.json",
	)
	// Output:
	//
	// false
	// "D3830"        "2022-07-20 18:21"
	//
	// true
}

func ExampleKunming_VehicleNo() {
	printVehicleNo(
		&adapters.Kunming{},
		"kunming_full.json",
		"kunming_invalid.json",
	)
	// Output:
	// "CRH2A4088"    false
	// ""             false
}
