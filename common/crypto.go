package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"

	"github.com/andreburgaud/crypt2go/ecb"
)

// DesEcbEncrypt encrypts the plain text with PKCS #7 padding and
// electronic codebook mode of operation.
func DesEcbEncrypt(plainText, key []byte) (cipherText []byte) {
	block, err := des.NewCipher(key)
	Must(err)
	plainText = PKCS7Padding(plainText, len(key))
	cipherText = make([]byte, len(plainText))
	blockMode := ecb.NewECBEncrypter(block)
	blockMode.CryptBlocks(cipherText, plainText)
	return
}

// DesEcbDecrypt is the counterpart of DesEcbEncrypt; it decrypts the cipher
// text and strips the PKCS #7 padding bytes off the end of the plain text.
func DesEcbDecrypt(cipherText, key []byte) (plainText []byte) {
	block, err := des.NewCipher(key)
	Must(err)
	plainText = make([]byte, len(cipherText))
	blockMode := ecb.NewECBDecrypter(block)
	blockMode.CryptBlocks(plainText, cipherText)
	return PKCS7Unpadding(plainText)
}

// AesEcbEncrypt encrypts the plain text with PKCS #7 padding and
// electronic codebook mode of operation.
func AesEcbEncrypt(plainText, key []byte) (cipherText []byte) {
	block, err := aes.NewCipher(key)
	Must(err)
	plainText = PKCS7Padding(plainText, len(key))
	cipherText = make([]byte, len(plainText))
	blockMode := ecb.NewECBEncrypter(block)
	blockMode.CryptBlocks(cipherText, plainText)
	return
}

// AesEcbDecrypt is the counterpart of AesEcbEncrypt; it decrypts the cipher
// text and strips the PKCS #7 padding bytes off the end of the plain text.
func AesEcbDecrypt(cipherText, key []byte) (plainText []byte) {
	block, err := aes.NewCipher(key)
	Must(err)
	plainText = make([]byte, len(cipherText))
	blockMode := ecb.NewECBDecrypter(block)
	blockMode.CryptBlocks(plainText, cipherText)
	return PKCS7Unpadding(plainText)
}

// AesCbcEncrypt encrypts the plain text with PKCS #7 padding, block chaining
// mode of operation, and a predefined initial vector.
func AesCbcEncrypt(plainText, key, iv []byte) (cipherText []byte) {
	block, err := aes.NewCipher(key)
	Must(err)
	plainText = PKCS7Padding(plainText, len(iv))
	cipherText = make([]byte, len(plainText))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherText, plainText)
	return
}

// AesCbcDecrypt is the counterpart of AesCbcEncrypt; it decrypts the cipher
// text and strips the PKCS #7 padding bytes off the end of the plain text.
func AesCbcDecrypt(cipherText, key, iv []byte) (plainText []byte) {
	block, err := aes.NewCipher(key)
	Must(err)
	plainText = make([]byte, len(cipherText))
	blockMode := cipher.NewCBCDecrypter(block, iv)
	blockMode.CryptBlocks(plainText, cipherText)
	return PKCS7Unpadding(plainText)
}

// PKCS7Padding pads the input octet vector to a multiple of blockSize octets
// with the scheme defined in RFC 2315.
func PKCS7Padding(input []byte, blockSize int) (buf []byte) {
	if len(input) == 0 || blockSize < 1 || blockSize > 255 {
		return
	}
	pad := blockSize - len(input)%blockSize
	buf = make([]byte, len(input)+pad)
	copy(buf, input)
	copy(buf[len(input):], bytes.Repeat([]byte{byte(pad)}, pad))
	return
}

// PKCS7Unpadding removes the padded bytes from the decrypted text
// according to the last decrypted byte to recover the original payload.
func PKCS7Unpadding(padded []byte) []byte {
	length := len(padded)
	if length == 0 {
		return nil
	}
	return padded[:length-int(padded[length-1])]
}
