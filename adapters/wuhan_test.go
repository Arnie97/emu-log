package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleWuhan_TrainNo() {
	printTrainNo(adapters.Wuhan{}, "wuhan_full.html")
	printTrainNo(adapters.Wuhan{}, "wuhan_basic.html")
	printTrainNo(adapters.Wuhan{}, "wuhan_invalid.html")

	// Output:
	// "G6856"        false ""
	// ""             true  ""
	// ""             true  ""
}

func ExampleWuhan_VehicleNo() {
	printVehicleNo(adapters.Wuhan{}, "wuhan_full.html")
	printVehicleNo(adapters.Wuhan{}, "wuhan_basic.html")
	printVehicleNo(adapters.Wuhan{}, "wuhan_invalid.html")

	// Output:
	// "CRH380A2900"  false
	// "CRH2A2445"    false
	// ""             true
}
