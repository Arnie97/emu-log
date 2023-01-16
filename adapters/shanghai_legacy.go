package adapters

import (
	"fmt"
	"net/http"

	"github.com/arnie97/emu-log/common"
)

type ShanghaiLegacy struct{}

func init() {
	Register(ShanghaiLegacy{})
}

func (ShanghaiLegacy) Code() string {
	return "H"
}

func (ShanghaiLegacy) Name() string {
	return "华东印记（咻咻365）"
}

func (ShanghaiLegacy) URL() (pattern string, mockValue interface{}) {
	return "http://portal.xiuxiu365.cn/portal/qrcode/%s", nil
}

func (ShanghaiLegacy) AlwaysOn() bool {
	return true
}

func (a ShanghaiLegacy) Info(serial string) (info JSONObject, err error) {
	const api = "https://g.xiuxiu365.cn/railway_api/web/index/train?pqCode=%s"
	url := fmt.Sprintf(api, serial)

	var resp *http.Response
	if resp, err = httpClient(a).Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status int `json:"code"`
		Msg    string
		Data   JSONObject
	}
	err = parseResult(resp, &result)
	info = result.Data
	return
}

func (ShanghaiLegacy) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	train := TrainSchedule{
		TrainNo: info["trainName"].(string),
	}
	if len(train.TrainNo) != 0 {
		trains = []TrainSchedule{train}
	}
	return
}

func (ShanghaiLegacy) UnitNo(_ string, info JSONObject) (unitNo string, err error) {
	defer common.Catch(&err)
	unitNo = common.NormalizeUnitNo(info["cdh"].(string))
	return
}
