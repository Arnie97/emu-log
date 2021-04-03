package adapters

import (
	"fmt"
	"net/http"

	"github.com/arnie97/emu-log/common"
)

type Shanghai struct{}

func init() {
	Register(Shanghai{})
}

func (Shanghai) Code() string {
	return "H"
}

func (Shanghai) Name() string {
	return "中国铁路上海局集团有限公司"
}

func (Shanghai) URL() string {
	return "http://portal.xiuxiu365.cn/portal/qrcode/%s"
}

func (Shanghai) BruteForce(pqCodes chan<- string) {
	for i := 2000; i < 11000; i += 200 {
		pqCodes <- fmt.Sprintf("PQ%07d", i)
	}
	for i := 11000; i < 2500000; i += 500 {
		pqCodes <- fmt.Sprintf("PQ%07d", i)
	}
}

func (Shanghai) AlwaysOn() bool {
	return true
}

func (Shanghai) Info(serial string) (info JSONObject, err error) {
	const api = "https://g.xiuxiu365.cn/railway_api/web/index/train?pqCode=%s"
	url := fmt.Sprintf(api, serial)

	var resp *http.Response
	if resp, err = common.HTTPClient().Get(url); err != nil {
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

func (Shanghai) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	train := TrainSchedule{
		TrainNo: info["trainName"].(string),
	}
	if len(train.TrainNo) != 0 {
		trains = []TrainSchedule{train}
	}
	return
}

func (Shanghai) VehicleNo(info JSONObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	vehicleNo = common.NormalizeVehicleNo(info["cdh"].(string))
	return
}
