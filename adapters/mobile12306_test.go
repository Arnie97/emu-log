package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleMobile12306_TrainNo() {
	printTrainNo(
		&adapters.Mobile12306{},
		"mobile12306_full.json",
		"mobile12306_invalid.json",
	)
	// Output:
	//
	// false
	// "D3830"        "2022-07-20 18:21"
	//
	// true
}

func ExampleMobile12306_UnitNo() {
	printUnitNo(
		&adapters.Mobile12306{},
		"mobile12306_full.json",
		"mobile12306_invalid.json",
	)
	// Output:
	// M "CRH2A4088"    false false false
	// H ""              true  true false
}
