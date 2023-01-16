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

type JinanQuery struct {
	Params    string `json:"params"`
	Timestamp int64  `json:"timeStamp"`
	CGUID     string `json:"cguid"`
	Token     string `json:"token,omitempty"`
	IsSign    int    `json:"isSign"`
	Signature string `json:"sign,omitempty"`
}

type Jinan struct{}

func init() {
	Register(Jinan{})
}

func (Jinan) Code() string {
	return "K"
}

func (Jinan) Name() string {
	return "掌上济铁（国铁吉讯）"
}

func (Jinan) URL() (pattern string, mockValue interface{}) {
	return "https://static.ccrgt.com/orderMeals?scene=%s", nil
}

func (Jinan) AlwaysOn() bool {
	return true
}

func (a Jinan) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("referer", fmt.Sprintf(
		"https://servicewechat.com/%s/54/page-frame.html", jinanApp,
	))
	return AdapterConf(a).Request.RoundTrip(req)
}

func (a Jinan) Info(serial string) (info JSONObject, err error) {
	return a.EncryptedQuery(
		"https://apicloud.ccrgt.com/crgt/retail-takeout/h5/takeout/scan/list",
		struct {
			SeatCode string `json:"seatCode"`
		}{serial},
	)
}

func (a Jinan) EncryptedQuery(api string, params interface{}) (info JSONObject, err error) {
	query := JinanQuery{
		Params:    a.InfoEncrypt(params),
		Timestamp: common.UnixMilli(),
		Token:     SessionID(a),
		IsSign:    2,
	}
	query.Signature = query.Sign()

	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(query); err != nil {
		return
	}
	buf := bytes.NewBuffer(jsonBytes)

	var resp *http.Response
	if resp, err = httpClient(a).Post(api, common.ContentType, buf); err != nil {
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
	err = a.InfoDecrypt(result.Data, &info)
	return
}

// InfoEncrypt encrypts the JSON string in AES-CBC-128 cipher mode, and return
// the cipher text encoded with padded base64 encoding scheme.
func (Jinan) InfoEncrypt(src interface{}) string {
	plainText, err := json.Marshal(src)
	common.Must(err)

	cipherText := common.AesCbcEncrypt(plainText, jinanKey, jinanIV)
	return base64.StdEncoding.EncodeToString(cipherText)
}

// InfoDecrypt decrypts the base64 encoded cipher text with AES-CBC-128,
// and unmarshals the plain text result into the given structure.
func (Jinan) InfoDecrypt(src string, dest interface{}) (err error) {
	cipherText, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return
	}

	plainText := common.AesCbcDecrypt(cipherText, jinanKey, jinanIV)
	return json.Unmarshal(plainText, dest)
}

// Sign serializes the message in a deterministic manner,
// and generates its hexadecimal encoded MD5 digest.
func (q JinanQuery) Sign() string {
	message := fmt.Sprintf(
		"%s%s%v%s%s",
		jinanApp,
		q.Token,
		q.Timestamp,
		base64.StdEncoding.EncodeToString(jinanKey),
		q.Params,
	)
	hash := md5.Sum([]byte(message))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

func (Jinan) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
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
			timestamp := elem.TrainInfo.ArriveTimestamp
			train.Date = time.Unix(timestamp, 0).Format(common.ISODate)
		} else {
			train.TrainNo += "/"
		}
		train.TrainNo += elem.TrainInfo.TrainNumber
	}
	trains = []TrainSchedule{train}
	return
}

func (Jinan) UnitNo(_ string, info JSONObject) (unitNo string, err error) {
	defer common.Catch(&err)
	unitNo = common.NormalizeUnitNo(info["czNo"].(string))
	return
}
