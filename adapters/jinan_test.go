package adapters_test

import (
	"fmt"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
)

func ExampleJinan_SerialEncrypt() {
	fmt.Println(adapters.Jinan{}.SerialEncrypt("K1001641584150"))
	// Output: 7x2dNPLNuBHTpL9Bc6z2JiKabX+sFyBLS4w1L0Ulbkw=
}

func ExampleJinan_AESEncrypt() {
	fmt.Println(adapters.Jinan{}.AESEncrypt(
		[]byte("Arnie97"),
		[]byte("$ecure*P@$$w0rd$"),
		[]byte("initialVector128"),
	))
	// Output: [46 169 15 51 223 19 237 171 243 81 115 177 56 118 214 219]
}

func ExampleJinan_PKCS7Padding() {
	fmt.Println(adapters.Jinan{}.PKCS7Padding([]byte("abcdefgh"), 8))
	fmt.Println(adapters.Jinan{}.PKCS7Padding([]byte("abcdefg"), 16))
	fmt.Println(adapters.Jinan{}.PKCS7Padding([]byte("abcdef"), 256))
	// Output:
	// [97 98 99 100 101 102 103 104 8 8 8 8 8 8 8 8]
	// [97 98 99 100 101 102 103 9 9 9 9 9 9 9 9 9]
	// []
}

func ExampleJinan_TrainNo() {
	common.SetMockHTTPClientRespBody(jinanFullResult)
	fmt.Println(adapters.Jinan{}.TrainNo(""))
	// Output: G2079/G2078 2020-07-09 <nil>
}

func ExampleJinan_VehicleNo() {
	common.SetMockHTTPClientRespBody(jinanFullResult)
	fmt.Println(adapters.Jinan{}.VehicleNo(""))
	// Output: CRH380B5847 <nil>
}
