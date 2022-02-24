package adapters

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/arnie97/emu-log/common"
)

var (
	harbinTrainNoRegExp = regexp.MustCompile(`<div class="cczi">([\w/]+)&nbsp;<font>([-\w/]*)</font></div>`)
)

type Harbin struct {
}

func init() {
	Register(Harbin{})
}

func (Harbin) Code() string {
	return "B"
}

func (Harbin) Name() string {
	return "中国铁路哈尔滨局集团有限公司"
}

func (Harbin) URL() (pattern string, mockValue interface{}) {
	return "http://l.jeehon.com/lkfw/api?id=%s", nil
}

func (Harbin) BruteForce(serials chan<- string) {
}

func (Harbin) AlwaysOn() bool {
	return false
}

func (b Harbin) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Cookie", common.Conf(b.Code()))
	return common.IntervalTransport{}.RoundTrip(req)
}

func (b Harbin) Info(serial string) (info JSONObject, err error) {
	const api = "http://l.jeehon.com/lkfw/api/index.asp?id=%s"
	url := fmt.Sprintf(api, serial)
	var resp *http.Response
	if resp, err = common.HTTPClient(b).Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	var bytes []byte
	bytes, err = ioutil.ReadAll(resp.Body)
	if match := harbinTrainNoRegExp.FindSubmatch(bytes); match != nil {
		info = JSONObject{
			"train": string(match[1]),
			"seat":  string(match[2]),
		}
	}
	return
}

func (b Harbin) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	trains = []TrainSchedule{{
		TrainNo: info["train"].(string),
	}}
	return
}

func (b Harbin) VehicleNo(serialNo string, info JSONObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	vehicleNo = fmt.Sprintf("CRH380BG5@%s", serialNo[:2])
	_, err = b.TrainNo(info)
	return
}
