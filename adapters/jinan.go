package adapters

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
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

func (Jinan) AlwaysOn() bool {
	return true
}

func (Jinan) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("referer", fmt.Sprintf(
		"https://servicewechat.com/%s/54/page-frame.html", jinanApp,
	))
	return common.IntervalTransport{}.RoundTrip(req)
}

func (b Jinan) Info(serial string) (info JSONObject, err error) {
	const api = "https://apicloud.ccrgt.com/crgt/retail-takeout/h5/takeout/scan/list"
	values := JSONObject{
		"params":    b.SerialEncrypt(serial),
		"timeStamp": time.Now().UnixNano() / 1000000,
		"cguid":     "",
		"token":     common.Conf(b.Code()),
		"isSign":    2,
	}
	values["sign"] = b.Signature(values)

	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(values); err != nil {
		return
	}
	buf := bytes.NewBuffer(jsonBytes)

	var resp *http.Response
	if resp, err = common.HTTPClient(b).Post(api, common.ContentType, buf); err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status int    `json:"code"`
		Msg    string `json:"errmsg"`
		Data   string `json:"data"`
	}
	if err = parseResult(resp, &result); err != nil {
		return
	}
	err = b.InfoDecrypt(result.Data, &info)
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

	cipherText := common.AESEncrypt(plainText, jinanKey, jinanIV)
	return base64.StdEncoding.EncodeToString(cipherText)
}

// InfoDecrypt decrypts the base64 encoded cipher text with AES-CBC-128,
// and unmarshals the plain text result into the given structure.
func (b Jinan) InfoDecrypt(src string, dest interface{}) (err error) {
	cipherText, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return
	}

	plainText := common.AESDecrypt(cipherText, jinanKey, jinanIV)
	return json.Unmarshal(plainText, dest)
}

// Signature serializes the message in a deterministic manner,
// and generates its hexadecimal encoded MD5 digest.
func (b Jinan) Signature(values JSONObject) string {
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

func (b Jinan) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
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
	if len(infoList) == 0 {
		return
	}

	var train TrainSchedule
	for i, elem := range infoList {
		if i == 0 {
			timestamp := elem.TrainInfo.StartTimestamp
			train.Date = time.Unix(timestamp, 0).Format(common.ISODate)
		} else {
			train.TrainNo += "/"
		}
		train.TrainNo += elem.TrainInfo.TrainNumber
	}
	trains = []TrainSchedule{train}
	return
}

func (Jinan) VehicleNo(info JSONObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	vehicleNo = common.NormalizeVehicleNo(info["czNo"].(string))
	return
}
