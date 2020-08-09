package adapters_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/arnie97/emu-log/adapters"
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

func ExampleAESEncrypt() {
	fmt.Println(adapters.AESEncrypt(
		[]byte("Arnie97"),
		[]byte("$ecure*P@$$w0rd$"),
		[]byte("initialVector128"),
	))
	// Output: [46 169 15 51 223 19 237 171 243 81 115 177 56 118 214 219]
}

func ExampleAESDecrypt() {
	entropy := make([]byte, 42)
	rand.Read(entropy)
	key, iv, text := entropy[:16], entropy[16:32], entropy[32:]
	cipherText := adapters.AESEncrypt(text, key, iv)
	fmt.Println(bytes.Compare(adapters.AESDecrypt(cipherText, key, iv), text))
	// Output:
	// 0
}

func ExamplePKCS7Padding() {
	fmt.Println(adapters.PKCS7Padding([]byte("abcdefgh"), 8))
	fmt.Println(adapters.PKCS7Padding([]byte("abcdefg"), 16))
	fmt.Println(adapters.PKCS7Padding([]byte("abcdef"), 256))
	// Output:
	// [97 98 99 100 101 102 103 104 8 8 8 8 8 8 8 8]
	// [97 98 99 100 101 102 103 9 9 9 9 9 9 9 9 9]
	// []
}

func ExamplePKCS7Unpadding() {
	fmt.Println(adapters.PKCS7Unpadding([]byte{}))
	fmt.Println(adapters.PKCS7Unpadding([]byte{1, 2, 3, 5, 5, 5, 5, 5}))
	// Output:
	// []
	// [1 2 3]
}

func ExampleJinan_Signature() {
	var exampleInput map[string]interface{}
	json.Unmarshal(readMockFile("jinan_input.json"), &exampleInput)
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
	// "G297/G300"    false "2020-08-10"
	// ""             false ""
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
