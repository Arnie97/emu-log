package common_test

import (
	"bytes"
	"fmt"
	"math/rand"

	"github.com/arnie97/emu-log/common"
)

func ExampleDesEcbEncrypt() {
	fmt.Println(common.DesEcbEncrypt(
		[]byte("Arnie97"),
		[]byte("P@$$w0rd"),
	))
	// Output: [175 255 31 191 150 239 19 134]
}

func ExampleDesEcbDecrypt() {
	entropy := make([]byte, 100)
	rand.Read(entropy)
	key, text := entropy[:8], entropy[8:]
	cipherText := common.DesEcbEncrypt(text, key)
	fmt.Println(bytes.Compare(common.DesEcbDecrypt(cipherText, key), text))
	// Output:
	// 0
}

func ExampleAesEcbEncrypt() {
	fmt.Println(common.AesEcbEncrypt(
		[]byte("Arnie97"),
		[]byte("$ecure*P@$$w0rd$"),
	))
	// Output: [65 104 5 194 136 199 209 10 240 48 109 82 11 205 74 176]
}

func ExampleAesEcbDecrypt() {
	entropy := make([]byte, 100)
	rand.Read(entropy)
	key, text := entropy[:32], entropy[32:]
	cipherText := common.AesEcbEncrypt(text, key)
	fmt.Println(bytes.Compare(common.AesEcbDecrypt(cipherText, key), text))
	// Output:
	// 0
}

func ExampleAesCbcEncrypt() {
	fmt.Println(common.AesCbcEncrypt(
		[]byte("Arnie97"),
		[]byte("$ecure*P@$$w0rd$"),
		[]byte("initialVector128"),
	))
	// Output: [46 169 15 51 223 19 237 171 243 81 115 177 56 118 214 219]
}

func ExampleAesCbcDecrypt() {
	entropy := make([]byte, 42)
	rand.Read(entropy)
	key, iv, text := entropy[:16], entropy[16:32], entropy[32:]
	cipherText := common.AesCbcEncrypt(text, key, iv)
	fmt.Println(bytes.Compare(common.AesCbcDecrypt(cipherText, key, iv), text))
	// Output:
	// 0
}

func ExamplePKCS7Padding() {
	fmt.Println(common.PKCS7Padding([]byte("abcdefgh"), 8))
	fmt.Println(common.PKCS7Padding([]byte("abcdefg"), 16))
	fmt.Println(common.PKCS7Padding([]byte("abcdef"), 256))
	// Output:
	// [97 98 99 100 101 102 103 104 8 8 8 8 8 8 8 8]
	// [97 98 99 100 101 102 103 9 9 9 9 9 9 9 9 9]
	// []
}

func ExamplePKCS7Unpadding() {
	fmt.Println(common.PKCS7Unpadding([]byte{}))
	fmt.Println(common.PKCS7Unpadding([]byte{1, 2, 3, 5, 5, 5, 5, 5}))
	// Output:
	// []
	// [1 2 3]
}
