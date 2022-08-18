package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleJD12306_TrainNo() {
	printTrainNo(
		&adapters.JD12306{},
		"jd12306_full.json",
		"jd12306_basic.json",
		"jd12306_invalid.json",
	)
	// Output:
	//
	// false
	// "D2751"        "2020-07-12"
	//
	// true
	//
	// true
}

func ExampleJD12306_UnitNo() {
	printUnitNo(
		&adapters.JD12306{},
		"jd12306_full.json",
		"jd12306_basic.json",
		"jd12306_invalid.json",
	)
	// Output:
	// "CRH5G5194"    false
	// "CRH380B3667"  false
	// ""             true
}
