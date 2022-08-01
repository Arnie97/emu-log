package adapters

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/arnie97/emu-log/common"
)

type Kunming struct {
}

func init() {
	Register(&Kunming{})
}

func (Kunming) Code() string {
	return "M"
}

func (Kunming) Name() string {
	return "中国铁路昆明局集团有限公司"
}

func (Kunming) URL() (pattern string, mockValue interface{}) {
	return "https://p.12306.cn/tservice/qr/travel/v1?c=%s-%02d-%02d%v&w=h", "F"
}

func (Kunming) AlwaysOn() bool {
	return false
}

func (b *Kunming) Info(serial string) (info JSONObject, err error) {
	const api = "https://mobile.12306.cn/wxxcx/wechat/main/travelServiceDecodeQrcode"

	urlStruct, err := url.ParseRequestURI(BuildURL(b, serial))
	if err != nil {
		return
	}
	form := urlStruct.Query()

	var resp *http.Response
	if resp, err = common.HTTPClient().PostForm(api, form); err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status bool       `json:"status"`
		Code   string     `json:"errorCode"`
		Msg    string     `json:"errorMsg"`
		Data   JSONObject `json:"data"`
	}
	err = parseResult(resp, &result)
	info = result.Data
	return
}

func (Kunming) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	shortDate := info["endDay"].(string)
	shortTime := info["endTime"].(string)
	trains = []TrainSchedule{{
		TrainNo: info["trainCode"].(string),
		Date: fmt.Sprintf(
			"%s-%s-%s %s:%s",
			shortDate[:4], shortDate[4:6], shortDate[6:8],
			shortTime[:2], shortTime[2:],
		),
	}}
	return
}

func (Kunming) VehicleNo(serialNo string, info JSONObject) (vehicleNo string, err error) {
	retrievedVehicleNo, _ := info["carCode"].(string)
	if len(retrievedVehicleNo) == 0 {
		// pass
	} else if serialNo != retrievedVehicleNo {
		vehicleNo = "CRH@" + retrievedVehicleNo
	} else {
		vehicleNo = common.NormalizeVehicleNo(retrievedVehicleNo)
	}
	return
}
