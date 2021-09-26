package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleShanghaiLegacy_BruteForce() {
	assertBruteForceRegExp(adapters.ShanghaiLegacy{}, `^PQ\d{7}$`)
	// Output:
}

func ExampleShanghaiLegacy_TrainNo() {
	printTrainNo(
		adapters.ShanghaiLegacy{},
		"shanghai_legacy_full.json",
		"shanghai_legacy_basic.json",
		"shanghai_legacy_invalid.json",
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

func ExampleShanghaiLegacy_VehicleNo() {
	printVehicleNo(adapters.ShanghaiLegacy{},
		"shanghai_legacy_full.json",
		"shanghai_legacy_basic.json",
		"shanghai_legacy_invalid.json",
	)
	// Output:
	// "CRH2A2001"    false
	// "CRH2C2150"    false
	// ""             true
}
