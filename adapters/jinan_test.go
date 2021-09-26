package adapters_test

import (
	"encoding/json"
	"fmt"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
)

func ExampleJinan_SerialEncrypt() {
	fmt.Println(adapters.Jinan{}.SerialEncrypt("K1001641584150"))
	// Output: 7x2dNPLNuBHTpL9Bc6z2JiKabX+sFyBLS4w1L0Ulbkw=
}

func ExampleJinan_InfoDecrypt() {
	var info interface{}
	adapters.Jinan{}.InfoDecrypt(
		"9RG2W94kpJJa9wxaAWL1/849sBfzHmbkmO7X3fBV1DU=", &info,
	)
	fmt.Println(info)
	// Output:
	// map[czNo:CRH380BL-5501]
}

func ExampleJinan_Signature() {
	var exampleInput map[string]interface{}
	json.Unmarshal(common.ReadMockFile("jinan_input.json"), &exampleInput)
	fmt.Println(adapters.Jinan{}.Signature(exampleInput) == exampleInput["sign"])
	// Output:
	// true
}

func ExampleJinan_BruteForce() {
	assertBruteForceRegExp(adapters.Jinan{}, `^K1001\d{9}$`)
	// Output:
}

func ExampleJinan_TrainNo() {
	printTrainNo(
		adapters.Jinan{},
		"jinan_full.json",
		"jinan_basic.json",
	)
	// Output:
	//
	// false
	// "G297/G300"    "2020-08-10"
	//
	// false
}

func ExampleJinan_VehicleNo() {
	printVehicleNo(
		adapters.Jinan{},
		"jinan_full.json",
		"jinan_basic.json",
	)
	// Output:
	// "CR400AFA2115" false
	// "CRH380BL5501" false
}
