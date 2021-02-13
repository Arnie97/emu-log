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
