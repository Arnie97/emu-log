package adapters

import (
	"encoding/json"
	"net/http"

	"github.com/arnie97/emu-log/common"
)

type Zhengzhou struct{}

func init() {
	Register(Zhengzhou{})
}

func (Zhengzhou) Code() string {
	return "F"
}

func (Zhengzhou) Name() string {
	return "中国铁路郑州局集团有限公司"
}

func (Zhengzhou) URL() string {
	return "https://p.12306.cn/tservice/catering/init?c=%s&w=h"
}

func (Zhengzhou) BruteForce(serials chan<- string) {
}

func (b Zhengzhou) Info(serial string) (info jsonObject, err error) {
	const api = "https://p.12306.cn/tservice/mealAction/qrcodeDecode"
	req, err := http.NewRequest("POST", api, nil)
	if err != nil {
		return
	}
	req.Header.Set("Cookie", common.Conf(b.Code()))
	resp, err := common.HTTPClient().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Data struct {
			TrainQrcodeInfo jsonObject
		}
	}
	err = json.NewDecoder(resp.Body).Decode(&result)
	info = result.Data.TrainQrcodeInfo
	return
}

func (b Zhengzhou) TrainNo(info jsonObject) (trainNo, date string, err error) {
	defer common.Catch(&err)
	trainNo = info["trainCode"].(string)
	date = info["startDay"].(string)
	date = date[:4] + "-" + date[4:6] + "-" + date[6:8]
	return
}

func (b Zhengzhou) VehicleNo(info jsonObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	vehicleNo = common.NormalizeVehicleNo(info["carCode"].(string))
	return
}
