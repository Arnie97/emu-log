package adapters

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/arnie97/emu-log/common"
)

var (
	harbinTrainNoRegExp = regexp.MustCompile(`<div class="cczi">(\w+)`)
)

type Harbin struct {
}

func init() {
	Register(&Harbin{})
}

func (Harbin) Code() string {
	return "B"
}

func (Harbin) Name() string {
	return "中国铁路哈尔滨局集团有限公司"
}

func (Harbin) URL() string {
	return "http://l.jeehon.com/lkfw/api?id=%s"
}

func (Harbin) BruteForce(serials chan<- string) {
}

func (Harbin) AlwaysOn() bool {
	return true
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
		info = JSONObject{b.Code(): string(match[1])}
	}
	return
}

func (b Harbin) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	trains = []TrainSchedule{{
		TrainNo: info[b.Code()].(string),
		Date:    time.Now().Format(common.ISODate),
	}}
	return
}

func (Harbin) VehicleNo(info JSONObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	return
}
