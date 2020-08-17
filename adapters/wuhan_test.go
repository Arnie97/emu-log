package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleWuhan_BruteForce() {
	assertBruteForceRegExp(adapters.Wuhan{}, `^\d{3}$`)
	// Output:
}

func ExampleWuhan_TrainNo() {
	printTrainNo(
		adapters.Wuhan{},
		"wuhan_full.html",
		"wuhan_anonymous.html",
		"wuhan_basic.html",
		"wuhan_invalid.html",
	)
	// Output:
	// "G1730"        false ""
	// "G551"         false ""
	// ""             true  ""
	// ""             true  ""
}

func ExampleWuhan_VehicleNo() {
	printVehicleNo(
		adapters.Wuhan{},
		"wuhan_full.html",
		"wuhan_anonymous.html",
		"wuhan_basic.html",
		"wuhan_invalid.html",
	)
	// Output:
	// "CR400AF2151"  false
	// "CR400AF2158"  false
	// "CRH380D1545"  false
	// ""             true
}
