package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleLTDD_TrainNo() {
	printTrainNo(
		adapters.LTDD{},
		"ltdd_full.html",
		"ltdd_anonymous.html",
		"ltdd_basic.html",
		"ltdd_invalid.html",
	)
	// Output:
	//
	// false
	// "G1730"        ""
	//
	// false
	// "G551"         ""
	//
	// true
	//
	// true
}

func ExampleLTDD_UnitNo() {
	printUnitNo(
		adapters.LTDD{},
		"ltdd_full.html",
		"ltdd_anonymous.html",
		"ltdd_basic.html",
		"ltdd_invalid.html",
	)
	// Output:
	// "CR400AF2151"  false
	// "CR400AF2158"  false
	// "CRH380D1545"  false
	// ""             true
}
