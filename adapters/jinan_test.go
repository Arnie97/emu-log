package adapters_test

import (
	"encoding/json"
	"fmt"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
)

func ExampleJinan_InfoEncrypt() {
	fmt.Println(adapters.Jinan{}.InfoEncrypt(map[string]string{
		"seatCode": "K1001641584150",
	}))
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

func ExampleJinanQuery_Sign() {
	var exampleInput adapters.JinanQuery
	json.Unmarshal(common.ReadMockFile("jinan_input.json"), &exampleInput)
	fmt.Println(exampleInput.Sign() == exampleInput.Signature)
	// Output:
	// true
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

func ExampleJinan_UnitNo() {
	printUnitNo(
		adapters.Jinan{},
		"jinan_full.json",
		"jinan_basic.json",
	)
	// Output:
	// K "CR400AFA2115" false false false
	// K "CRH380BL5501" false false false
}
