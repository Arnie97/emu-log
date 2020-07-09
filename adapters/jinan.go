package adapters

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/arnie97/emu-log/common"
)

var (
	jinanKey = []byte("prod_CrgtKey2019")
	jinanIV  = []byte("prod_iv20191001H")
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

func (Jinan) BruteForce(serials chan<- string) {
}

func (b Jinan) Info(serial string) (info jsonObject, err error) {
	const api = "https://apicloud.ccrgt.com/crgt/retail-takeout/h5/takeout/scan/list"
	values := jsonObject{
		"params": b.SerialEncrypt(serial),
		"token":  common.Conf(b.Code()),
		"isSign": 1,
	}

	jsonBytes, err := json.Marshal(values)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(jsonBytes)
	resp, err := common.HTTPClient().Post(api, contentType, buf)
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

func (b Jinan) TrainNo(serial string) (trainNo, date string, err error) {
	var (
		info     jsonObject
		infoList []struct {
			TrainInfo struct {
				TrainNumber, FirstStation, LastStation string
				StartTimestamp, ArriveTimestamp        int64
			}
		}
	)
	if info, err = b.Info(serial); err != nil {
		return
	}
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

func (b Jinan) VehicleNo(serial string) (vehicleNo string, err error) {
	var info jsonObject
	info, err = b.Info(serial)
	if err == nil {
		defer common.Catch(&err)
		vehicleNo = common.NormalizeVehicleNo(info["czNo"].(string))
	}
	return
}
