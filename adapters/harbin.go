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

type Harbin struct{}

func init() {
	Register(Harbin{})
}

func (Harbin) Code() string {
	return "B"
}

func (Harbin) Name() string {
	return "旅客服务系统（继峰科技）"
}

func (Harbin) URL() (pattern string, mockValue interface{}) {
	return "http://l.jeehon.com/lkfw/api?id=%s", nil
}

func (Harbin) AlwaysOn() bool {
	return false
}

func (a Harbin) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Cookie", SessionID(a))
	return AdapterConf(a).Request.RoundTrip(req)
}

func (a Harbin) Info(serial string) (info JSONObject, err error) {
	const api = "http://l.jeehon.com/lkfw/api/index.asp?id=%s"
	url := fmt.Sprintf(api, serial)
	var resp *http.Response
	if resp, err = httpClient(a).Get(url); err != nil {
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

func (a Harbin) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	trains = []TrainSchedule{{
		TrainNo: info["train"].(string),
	}}
	return
}

func (a Harbin) UnitNo(serialNo string, info JSONObject) (unitNo string, err error) {
	defer common.Catch(&err)
	unitNo = fmt.Sprintf("CRH380BG5@%s", serialNo[:2])
	_, err = a.TrainNo(info)
	return
}
