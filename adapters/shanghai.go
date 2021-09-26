package adapters

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/arnie97/emu-log/common"
)

type Shanghai struct{}

const (
	shanghaiKey = "8ab0aa3e08b9ca4c"
)

func init() {
	Register(Shanghai{})
}

func (Shanghai) Code() string {
	return "U"
}

func (Shanghai) Name() string {
	return "中国铁路上海局集团有限公司（爱上铁）"
}

func (Shanghai) URL() string {
	return "https://ky.railshj.cn?CHN=orderfood&type=ARMRSET&qrCode=%s"
}

func (Shanghai) BruteForce(pqCodes chan<- string) {
}

func (Shanghai) AlwaysOn() bool {
	return false
}

func (b Shanghai) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("channel", "MALL-WX-APPLET")
	req.Header.Set("version", "MALL-WX-APPLET_1.0.7")
	return common.IntervalTransport{}.RoundTrip(req)
}

func (b Shanghai) Info(serial string) (info JSONObject, err error) {
	const api = "https://ky.railshj.cn/12306app/orderingfood/pqcode/getTrainByPqCode"
	buf := bytes.NewBuffer(b.SerialEncrypt(serial))

	var resp *http.Response
	if resp, err = common.HTTPClient(b).Post(api, common.ContentType, buf); err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Code   string `json:"returnCode"`
		Msg    string `json:"returnMsg"`
		Data   string `json:"data"`
		Status bool   `json:"success"`
	}
	err = parseResult(resp, &result)
	return
}

// SerialEncrypt first wraps the serial number in JSON,
// then calculates the hash digest signature and put that in another JSON key,
// next encode the JSON string with padded base64 encoding scheme,
// after that encrypts the base64-encoded string with AES-ECB cipher mode,
// and finally wrap the result in JSON again.
func (b Shanghai) SerialEncrypt(serial string) []byte {
	message := struct {
		SignExt      string `json:"signext,omitempty"`
		TimestampExt int64  `json:"timestampext,omitempty"`
		PQCode       string `json:"pqCode"`
	}{
		PQCode: serial,
	}
	message.SignExt = b.Signature(message)
	message.TimestampExt = common.UnixMilli()

	plainText, err := json.Marshal(message)
	common.Must(err)
	base64Str := base64.StdEncoding.EncodeToString(plainText)
	cipherText := common.AesEcbEncrypt([]byte(base64Str), []byte(shanghaiKey))

	wrapper := struct {
		Data []byte `json:"data"`
	}{
		Data: cipherText,
	}
	jsonBytes, err := json.Marshal(wrapper) // implies base64 encode
	common.Must(err)
	return jsonBytes
}

// Signature serializes the message, and generates its
// hexadecimal hash digest with the MD5 - SHA1 hash chain.
func (b Shanghai) Signature(message interface{}) string {
	jsonBytes, err := json.Marshal(message)
	common.Must(err)

	base64Str := base64.StdEncoding.EncodeToString(jsonBytes)

	md5Sum := md5.Sum([]byte(base64Str))
	md5SumHex := hex.EncodeToString(md5Sum[:])

	sha1Sum := sha1.Sum([]byte(md5SumHex))
	sha1SumHex := hex.EncodeToString(sha1Sum[:])

	return sha1SumHex[6:26]
}

func (Shanghai) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	return
}

func (Shanghai) VehicleNo(info JSONObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	return
}
