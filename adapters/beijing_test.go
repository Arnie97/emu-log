package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleBeijing_TrainNo() {
	printTrainNo(adapters.Beijing{}, "beijing_full.json")
	printTrainNo(adapters.Beijing{}, "beijing_invalid.json")

	// Output:
	// "G666"         false "2020-07-10"
	// ""             true  ""
}

func ExampleBeijing_VehicleNo() {
	printVehicleNo(adapters.Beijing{}, "beijing_full.json")
	printVehicleNo(adapters.Beijing{}, "beijing_invalid.json")

	// Output:
	// "CR400AF0207"  false
	// ""             true
}
