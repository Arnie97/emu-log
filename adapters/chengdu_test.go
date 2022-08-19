package adapters_test

import (
	"github.com/arnie97/emu-log/adapters"
)

func ExampleChengdu_TrainNo() {
	printTrainNo(
		adapters.Chengdu{},
		"chengdu_full_1.base64",
		"chengdu_full_2.base64",
		"chengdu_invalid.base64",
	)
	// Output:
	//
	// false
	// "C6206"        "2021-02-07 06:30"
	// "C6303"        "2021-02-07 08:09"
	// "C6304"        "2021-02-07 11:33"
	// "C6313"        "2021-02-07 15:03"
	// "C6336"        "2021-02-07 18:47"
	// "D6192"        "2021-02-08 07:10"
	// "C6307"        "2021-02-08 08:09"
	// "C6308"        "2021-02-08 12:43"
	// "C6315"        "2021-02-08 15:52"
	// "C6276"        "2021-02-08 19:36"
	//
	// false
	// "G2187"        "2021-02-05 09:36"
	// "G2185"        "2021-02-06 06:45"
	//
	// true
}

func ExampleChengdu_UnitNo() {
	printUnitNo(
		adapters.Chengdu{},
		"chengdu_full_1.base64",
		"chengdu_full_2.base64",
		"chengdu_invalid.base64",
	)
	// Output:
	// W "CRH3A3089"    false false false
	// W "CRH@1582"     false false false
	// W ""             false false false
}
