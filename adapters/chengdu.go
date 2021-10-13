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
	return "中国铁路成都局集团有限公司"
}

func (Chengdu) URL() string {
	return "https://kyd.cd-rail.com?code=%s"
}

func (Chengdu) BruteForce(serials chan<- string) {
}

func (Chengdu) AlwaysOn() bool {
	return false
}

func (b Chengdu) Info(serial string) (info JSONObject, err error) {
	const api = "https://kyd.cd-rail.com/KYDMS_S/WeixinServlet"

	var (
		vehicleNo string
		form      url.Values
	)
	if vehicleNo, form, err = b.SerialEncrypt(api, serial); err != nil {
		return
	}

	var resp *http.Response
	if resp, err = common.HTTPClient().PostForm(api, form); err != nil {
		return
	}
	defer resp.Body.Close()

	info = map[string]interface{}{b.URL(): vehicleNo}
	result := []*JSONObject{&info}
	if err = b.InfoDecrypt(resp.Body, &result); err != nil {
		return
	}
	return
}

// SerialEncrypt converts the QR code tuple to a form,
// and encrypts the form values in DES-ECB cipher mode.
func (b Chengdu) SerialEncrypt(api, serial string) (vehicleNo string, ret url.Values, err error) {
	components := strings.Split(serial, ",")
	if len(components) != 6 {
		err = fmt.Errorf("invalid serial number tuple: %s", serial)
		return
	}
	var (
		today         = time.Now().Format("2006/01/02")
		vehicleModel  = components[1]
		vehicleDigits = components[2]
		seatCoach     int
		sqlParams     []byte
	)
	if seatCoach, err = strconv.Atoi(components[3]); err != nil {
		return
	}

	// save the vehicle number in the QR code for later comparison
	vehicleNo = common.NormalizeVehicleNo(vehicleModel + "-" + vehicleDigits)

	sqlParamsTuple := []interface{}{vehicleDigits, seatCoach, today, vehicleModel}
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
func (b Chengdu) InfoDecrypt(src io.Reader, dest interface{}) (err error) {
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

func (b Chengdu) VehicleNo(info JSONObject) (vehicleNo string, err error) {
	retrievedVehicleNo, _ := info["TRAIN_UNDER"].(string)
	if len(retrievedVehicleNo) == 0 {
		return
	}

	vehicleNo, _ = info[b.URL()].(string)
	if common.ApproxEqualVehicleNo(vehicleNo, retrievedVehicleNo) {
		return
	}

	// the number in the QR code does not match the one in the HTTP response
	vehicleNo = "CRH@" + retrievedVehicleNo
	return
}
