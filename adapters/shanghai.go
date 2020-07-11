package adapters

import (
	"fmt"
	"net/url"

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
	for i := 11000; i < 1550000; i += 500 {
		pqCodes <- fmt.Sprintf("PQ%07d", i)
	}
}

func (Shanghai) Info(pqCode string) (info jsonObject, err error) {
	const api = "https://g.xiuxiu365.cn/railway_api/web/index/train"
	query := url.Values{"pqCode": {pqCode}}.Encode()
	resp, err := common.HTTPClient().Get(api + "?" + query)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status int `json:"code"`
		Msg    string
		Data   jsonObject
	}
	err = parseResult(resp, &result)
	info = result.Data
	return
}

func (b Shanghai) TrainNo(pqCode string) (trainNo, date string, err error) {
	var info jsonObject
	info, err = b.Info(pqCode)
	if err == nil {
		defer common.Catch(&err)
		trainNo = info["trainName"].(string)
	}
	return
}

func (b Shanghai) VehicleNo(pqCode string) (vehicleNo string, err error) {
	var info jsonObject
	info, err = b.Info(pqCode)
	if err == nil {
		defer common.Catch(&err)
		vehicleNo = common.NormalizeVehicleNo(info["cdh"].(string))
	}
	return
}
