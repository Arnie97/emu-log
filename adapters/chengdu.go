package adapters

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/arnie97/emu-log/common"
)

const (
	chengduKey   = "kyd@info"
	chengduCode  = "kyd_lyk0351"
	chengduLogin = `["10.192.111.79","hhs","hhs"]`
)

type Chengdu struct{}

func init() {
	Register(Chengdu{})
}

func (Chengdu) Code() string {
	return "W"
}

func (Chengdu) Name() string {
	return "川之味"
}

func (Chengdu) URL() (pattern string, mockValue interface{}) {
	return "https://kyd.cd-rail.com?code=%s", nil
}

func (Chengdu) AlwaysOn() bool {
	return false
}

func (a Chengdu) Info(serial string) (info JSONObject, err error) {
	const api = "https://kyd.cd-rail.com/KYDMS_S/WeixinServlet"

	var form url.Values
	if form, err = a.SerialEncrypt(api, serial); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = common.HTTPClient().PostForm(api, form); err != nil {
		return
	}
	defer resp.Body.Close()

	result := []*JSONObject{&info}
	if err = a.InfoDecrypt(resp.Body, &result); err != nil {
		return
	}
	return
}

// SerialEncrypt converts the QR code tuple to a form,
// and encrypts the form values in DES-ECB cipher mode.
func (Chengdu) SerialEncrypt(api, serial string) (ret url.Values, err error) {
	components := strings.Split(serial, ",")
	if len(components) != 6 {
		err = fmt.Errorf("invalid serial number tuple: %s", serial)
		return
	}
	var (
		today      = time.Now().Format("2006/01/02")
		unitClass  = components[1]
		unitSerial = components[2]
		seatCoach  int
		sqlParams  []byte
	)
	if seatCoach, err = strconv.Atoi(components[3]); err != nil {
		return
	}

	sqlParamsTuple := []interface{}{unitSerial, seatCoach, today, unitClass}
	if sqlParams, err = json.Marshal(sqlParamsTuple); err != nil {
		return
	}

	base64Encrypt := func(data string) string {
		return base64.StdEncoding.EncodeToString(
			common.DesEcbEncrypt([]byte(data), []byte(chengduKey)))
	}
	ret = url.Values{
		"code":  {base64Encrypt(chengduCode)},
		"sql":   {base64Encrypt(string(sqlParams))},
		"url":   {api},
		"type":  {"POST"},
		"where": {base64Encrypt("[]")},
		"order": {base64Encrypt("[]")},
		"login": {base64Encrypt(chengduLogin)},
	}
	return
}

// InfoDecrypt decrypts the base64 encoded cipher text with DES-ECB,
// and unmarshals the plain text result into the given structure.
func (Chengdu) InfoDecrypt(src io.Reader, dest interface{}) (err error) {
	defer common.Catch(&err)

	var (
		base64Encoded []byte
		cipherText    []byte
	)
	if base64Encoded, err = ioutil.ReadAll(src); err != nil {
		return
	}
	base64Str := string(base64Encoded)
	if cipherText, err = base64.StdEncoding.DecodeString(base64Str); err != nil {
		return
	}

	plainText := common.DesEcbDecrypt(cipherText, []byte(chengduKey))
	return json.Unmarshal(plainText, dest)
}

func (Chengdu) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)

	var (
		trainListBytes = []byte(info["TRAIN_ORDER_TIME"].(string))
		trainList      [][2]string
	)
	if err = json.Unmarshal(trainListBytes, &trainList); err != nil {
		return
	}
	for _, pair := range trainList {
		trains = append(trains, TrainSchedule{
			TrainNo: pair[0],
			Date:    strings.Replace(pair[1], "/", "-", 2),
		})
	}
	return
}

func (Chengdu) UnitNo(qrTuple string, info JSONObject) (unitNo string, err error) {
	retrievedUnitNo, _ := info["TRAIN_UNDER"].(string)
	if len(retrievedUnitNo) == 0 {
		return
	}

	var (
		components = strings.Split(qrTuple, ",")
		unitClass  = components[1]
		unitSerial = components[2]
	)
	unitNo = unitClass + unitSerial
	if common.ApproxEqualUnitNo(unitNo, retrievedUnitNo) {
		return
	}

	// the number in the QR code does not match the one in the HTTP response
	unitNo = "CRH@" + retrievedUnitNo
	return
}
