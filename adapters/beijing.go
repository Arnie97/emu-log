package adapters

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/url"

	"github.com/arnie97/emu-log/common"
)

type Beijing struct{}

func init() {
	Register(Beijing{})
}

func (Beijing) Code() string {
	return "P"
}

func (Beijing) Name() string {
	return "康之旅（京铁列服）"
}

func (Beijing) URL() (pattern string, mockValue interface{}) {
	return "https://aymaoto.jtlf.cn/page/oto/index?QR=%s", nil
}

var (
	turn  int
	shift = []int{0, 1000, 500, 1500}
)

func (Beijing) AlwaysOn() bool {
	return false
}

func (Beijing) Info(qrCode string) (info JSONObject, err error) {
	const api = "https://aymaoto.jtlf.cn/webapi/otoshopping/ewh_getqrcodetrainnoinfo"
	const key = "qrcode=%s&key=ltRsjkiM8IRbC80Ni1jzU5jiO6pJvbKd"
	sign := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf(key, qrCode))))
	form := url.Values{"qrCode": {qrCode}, "sign": {sign}}

	var resp *http.Response
	if resp, err = common.HTTPClient().PostForm(api, form); err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status int `json:"state"`
		Msg    string
		Data   struct {
			TrainInfo JSONObject
			URLStr    string
		}
	}
	err = parseResult(resp, &result)
	info = result.Data.TrainInfo
	return
}

func (Beijing) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	train := TrainSchedule{
		TrainNo: info["TrainnoId"].(string),
		Date:    info["TrainnoDate"].(string),
	}
	trains = []TrainSchedule{train}
	return
}

func (Beijing) UnitNo(_ string, info JSONObject) (unitNo string, err error) {
	defer common.Catch(&err)
	unitNo = common.NormalizeUnitNo(info["TrainId"].(string))
	return
}
