package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

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

func ExampleShanghaiLegacy_UnitNo() {
	printUnitNo(adapters.ShanghaiLegacy{},
		"shanghai_legacy_full.json",
		"shanghai_legacy_basic.json",
		"shanghai_legacy_invalid.json",
	)
	// Output:
	// H "CRH2A2001"    false false false
	// H "CRH2C2150"    false false false
	// H ""              true  true false
}
