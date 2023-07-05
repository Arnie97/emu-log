package adapters

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/arnie97/emu-log/common"
)

var (
	unitNoWithBureauCodeRegExp = regexp.MustCompile(`^[A-Z]\d{7}$`)
)

type Mobile12306 struct{}

func init() {
	Register(&Mobile12306{})
}

func (Mobile12306) Code() string {
	return "M"
}

func (Mobile12306) Name() string {
	return "铁路畅行"
}

func (Mobile12306) URL() (pattern string, mockValue interface{}) {
	return "https://p.12306.cn/tservice/qr/travel/v1?c=%s&w=%v", "h"
}

func (Mobile12306) AlwaysOn() bool {
	return false
}

func (a *Mobile12306) Info(serial string) (info JSONObject, err error) {
	const api = "https://mobile.12306.cn/wxxcx/wechat/main/travelServiceDecodeQrcode"

	urlStruct, err := url.ParseRequestURI(BuildURL(a, serial))
	if err != nil {
		return
	}
	form := urlStruct.Query()

	var resp *http.Response
	if resp, err = httpClient(a).PostForm(api, form); err != nil {
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

func (Mobile12306) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
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

func (Mobile12306) UnitNo(serialNo string, info JSONObject) (unitNo string, err error) {
	defer common.Catch(&err)
	unitNo = common.NormalizeUnitNo(info["carCode"].(string))
	return
}

func (Mobile12306) Operator(serialNo string, info JSONObject) (bureauCode string, err error) {
	if bureauCode, ok := info["bureauCode"].(string); ok {
		return bureauCode, nil
	}
	if unitNoWithBureauCodeRegExp.MatchString(serialNo) {
		return serialNo[:1], nil
	}
	return
}
