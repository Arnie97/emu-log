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

func ExampleNormalizeVehicleNo() {
	fmt.Println(common.NormalizeVehicleNo("CRH_6-002A"))
	// Output: CRH6002A
}
