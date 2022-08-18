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
	//
	// false
	// "G1363"        ""
	//
	// true
}

func ExampleGuangzhou_UnitNo() {
	printUnitNo(
		adapters.Guangzhou{},
		"guangzhou_full.json",
		"guangzhou_invalid.json",
	)
	// Output:
	// "CR16-11876@9" false
	// ""             true
}
