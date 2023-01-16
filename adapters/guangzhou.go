package adapters

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arnie97/emu-log/common"
)

type Guangzhou struct{}

func init() {
	Register(Guangzhou{})
}

func (Guangzhou) Code() string {
	return "Q"
}

func (Guangzhou) Name() string {
	return "舌尖上的旅途（易食纵横）"
}

func (Guangzhou) URL() (pattern string, mockValue interface{}) {
	return "https://sj-wake.yishizongheng.com/scanOrder?code=%s&carriage=%d&site=%v", "1F"
}

func (Guangzhou) AlwaysOn() bool {
	return false
}

func (a Guangzhou) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", SessionID(a))
	return AdapterConf(a).Request.RoundTrip(req)
}

func (a Guangzhou) Info(serial string) (info JSONObject, err error) {
	const api = "https://sj-api.yishizongheng.com/shejian/api/train/getByQrcode?qrcode=%s"
	url := fmt.Sprintf(api, strings.TrimLeft(serial, "0"))

	var resp *http.Response
	if resp, err = httpClient(a).Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status int    `json:"code"`
		Msg    string `json:"message"`
		Data   JSONObject
	}
	err = parseResult(resp, &result)
	if result.Data != nil {
		return result.Data, nil
	}
	return nil, err
}

func (Guangzhou) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	trains = []TrainSchedule{{
		TrainNo: info["train"].(string),
	}}
	return
}

func (Guangzhou) UnitNo(serialNo string, info JSONObject) (unitNo string, err error) {
	defer common.Catch(&err)
	unitNo = fmt.Sprintf(
		"CR%s-%.0f@%s",
		info["carriageNum"].(string), info["id"], serialNo,
	)
	return
}
