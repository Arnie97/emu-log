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
	// "D2936"        ""
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
	// Z "CRH2A2276"    false false false
	// N "CR400AF2158"  false false false
	// W "CRH380D1545"  false false false
	//   ""             false  true  true
}
