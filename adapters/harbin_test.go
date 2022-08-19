package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleHarbin_TrainNo() {
	printTrainNo(
		adapters.Harbin{},
		"harbin_basic.html",
		"harbin_invalid.html",
	)
	// Output:
	//
	// false
	// "G1206/7"      ""
	//
	// true
}

func ExampleHarbin_UnitNo() {
	printUnitNo(adapters.Harbin{},
		"harbin_basic.html",
		"harbin_invalid.html",
	)
	// Output:
	// B "CRH380BG5@42" false false false
	// B "CRH380BG5@42" false  true false
}
