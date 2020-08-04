package adapters

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/arnie97/emu-log/common"
)

var (
	jinanKey = []byte("prod_CrgtKey2019")
	jinanIV  = []byte("prod_iv20191001H")
	jinanApp = "wxc33f19505fa37f4e"
)

type Jinan struct{}

func init() {
	Register(Jinan{})
}

func (Jinan) Code() string {
	return "K"
}

func (Jinan) Name() string {
	return "中国铁路济南局集团有限公司"
}

func (Jinan) URL() string {
	return "https://static.ccrgt.com/orderMeals?scene=%s"
}

func (Jinan) BruteForce(serials chan<- string) {
}

func (b Jinan) Info(serial string) (info jsonObject, err error) {
	const api = "https://apicloud.ccrgt.com/crgt/retail-takeout/h5/takeout/scan/list"
	values := jsonObject{
		"params":    b.SerialEncrypt(serial),
		"timeStamp": time.Now().UnixNano() / 1000000,
		"cguid":     "",
		"token":     common.Conf(b.Code()),
		"isSign":    2,
	}
	values["sign"] = b.Signature(values)

	jsonBytes, err := json.Marshal(values)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(jsonBytes)
	resp, err := common.HTTPClient().Post(api, common.ContentType, buf)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status int    `json:"code"`
		Msg    string `json:"errmsg"`
		Data   jsonObject
	}
	err = parseResult(resp, &result)
	info = result.Data
	return
}

// SerialEncrypt first wraps the serial number in JSON, then encrypts
// the JSON string with AES-CBC-128 cipher mode, and finally return
// the cipher text encoded with padded base64 encoding scheme.
func (b Jinan) SerialEncrypt(serial string) string {
	plainText, err := json.Marshal(struct {
		SeatCode string `json:"seatCode"`
	}{serial})
	common.Must(err)

	cipherText := b.AESEncrypt(plainText, jinanKey, jinanIV)
	return base64.StdEncoding.EncodeToString(cipherText)
}

// AESEncrypt encrypts the plain text with PKCS #7 padding, block chaining
// mode of operation, and a predefined initial vector.
func (b Jinan) AESEncrypt(plainText, key, iv []byte) (cipherText []byte) {
	block, err := aes.NewCipher(key)
	common.Must(err)
	plainText = b.PKCS7Padding(plainText, len(iv))
	cipherText = make([]byte, len(plainText))
	blockMode := cipher.NewCBCEncrypter(block, iv)
	blockMode.CryptBlocks(cipherText, plainText)
	return
}

// PKCS7Padding pads the input octet vector to a multiple of blockSize octets
// with the scheme defined in RFC 2315.
func (b Jinan) PKCS7Padding(input []byte, blockSize int) (buf []byte) {
	if len(input) == 0 || blockSize < 1 || blockSize > 255 {
		return
	}
	pad := blockSize - len(input)%blockSize
	buf = make([]byte, len(input)+pad)
	copy(buf, input)
	copy(buf[len(input):], bytes.Repeat([]byte{byte(pad)}, pad))
	return
}

// Signature serializes the message in a deterministic manner,
// and generates its hexadecimal encoded MD5 digest.
func (b Jinan) Signature(values jsonObject) string {
	message := fmt.Sprintf(
		"%s%s%v%s%s",
		jinanApp,
		values["token"],
		values["timeStamp"],
		base64.StdEncoding.EncodeToString(jinanKey),
		values["params"],
	)
	hash := md5.Sum([]byte(message))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func (b Jinan) TrainNo(info jsonObject) (trainNo, date string, err error) {
	var (
		infoList []struct {
			TrainInfo struct {
				TrainNumber, FirstStation, LastStation string
				StartTimestamp, ArriveTimestamp        int64
			}
		}
	)
	if err = common.StructDecode(info["trainInfos"], &infoList); err != nil {
		return
	}

	for i, elem := range infoList {
		if i == 0 {
			timestamp := elem.TrainInfo.StartTimestamp
			date = time.Unix(timestamp, 0).Format(common.ISODate)
		} else {
			trainNo += "/"
		}
		trainNo += elem.TrainInfo.TrainNumber
	}
	return
}

func (b Jinan) VehicleNo(info jsonObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	vehicleNo = common.NormalizeVehicleNo(info["czNo"].(string))
	return
}
