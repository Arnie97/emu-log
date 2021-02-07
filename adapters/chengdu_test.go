package adapters_test

import (
	"bytes"
	"fmt"
	"math/rand"

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
	// "C6276"        false "2021-02-08 19:36"
	// "G2185"        false "2021-02-06 06:45"
	// ""             true  ""
}

func ExampleChengdu_VehicleNo() {
	printVehicleNo(
		adapters.Chengdu{},
		"chengdu_full_1.base64",
		"chengdu_full_2.base64",
		"chengdu_invalid.base64",
	)
	// Output:
	// "CRH3A3089"    false
	// "CRH@1582"     false
	// ""             false
}

func ExampleDESEncrypt() {
	fmt.Println(adapters.DESEncrypt(
		[]byte("Arnie97"),
		[]byte("P@$$w0rd"),
	))
	// Output: [175 255 31 191 150 239 19 134]
}

func ExampleDESDecrypt() {
	entropy := make([]byte, 100)
	rand.Read(entropy)
	key, text := entropy[:8], entropy[8:]
	cipherText := adapters.DESEncrypt(text, key)
	fmt.Println(bytes.Compare(adapters.DESDecrypt(cipherText, key), text))
	// Output:
	// 0
}
