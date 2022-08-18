package common_test

import (
	"fmt"

	"github.com/arnie97/emu-log/common"
)

func ExampleNormalizeTrainNo() {
	fmt.Println(
		common.NormalizeTrainNo("C1040/37/40"),
		common.NormalizeTrainNo("G1040/1"),
		common.NormalizeTrainNo("D1040"),
		common.NormalizeTrainNo("1040/1"),
		common.NormalizeTrainNo("CRH6"),
	)
	// Output: [C1040 C1037 C1040] [G1040 G1041] [D1040] [] []
}

func ExampleNormalizeUnitNo() {
	for _, unitNo := range []string{
		"CRH_6-002A",
		"CRH5A1-5028",
		"CR200J2-4001",
		"CHR380B-3770",
		"CRH380DV-1503",
		"CRH380D-V1-1504",
		"CR400BFB-1-5097 ",
		" CR400AFBZ2-2249 ",
	} {
		fmt.Println(common.NormalizeUnitNo(unitNo))
	}

	// Output:
	//
	// CRH6002A
	// CRH5A5028
	// CR200J4001
	// CRH380B3770
	// CRH380D1503
	// CRH380D1504
	// CR400BFB5097
	// CR400AFBZ2249
}

func ExampleApproxEqualUnitNo() {
	fmt.Println(
		common.ApproxEqualUnitNo("CRH380B3626", "CHR380B3626"),
		common.ApproxEqualUnitNo("CR400BF5033", "5033"),
		common.ApproxEqualUnitNo("CRH5A5124", "CRH5A15124"),
		common.ApproxEqualUnitNo("CRH2E2462", "CR8+8-0@459"),
		common.ApproxEqualUnitNo("CRH2A2002", "CRH6A4002"),
		common.ApproxEqualUnitNo("CR", "CR"),
		common.ApproxEqualUnitNo("CRH6C2145", ""),
		common.ApproxEqualUnitNo("", "CRH2C2150"),
	)
	// Output: true true true true false true false true
}
