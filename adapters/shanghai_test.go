package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleShanghai_BruteForce() {
	assertBruteForceRegExp(adapters.Shanghai{}, `^PQ\d{7}$`)
	// Output:
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
	// "D3074/D3071"  ""
	//
	// false
	//
	// true
}

func ExampleShanghai_VehicleNo() {
	printVehicleNo(adapters.Shanghai{},
		"shanghai_full.json",
		"shanghai_basic.json",
		"shanghai_invalid.json",
	)
	// Output:
	// "CRH2A2001"    false
	// "CRH2C2150"    false
	// ""             true
}
