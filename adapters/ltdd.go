package adapters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/arnie97/emu-log/common"
)

var (
	htmlUnitNoRegExp  = regexp.MustCompile(`<p class="schedule">.+:(\w+)-6-6D/F</p>`)
	htmlTrainNoRegExp = regexp.MustCompile(`<p class="train-title">(\w+)随车购</p>`)
	jsonUnitNoRegExp  = regexp.MustCompile(`var locomotive_info = (\{.+\});`)
	jsonCompanyRegExp = regexp.MustCompile(`var company_info = (\{.+\});`)
	jsonCompanyIDMap  = map[string]string{
		"4":  "Z", // 广西宁铁餐饮服务有限公司
		"25": "W", // 成都客运段
		"49": "N", // 新武汉动高餐饮管理服务有限公司
	}
)

type LTDD struct{}

func init() {
	Register(LTDD{})
}

func (LTDD) Code() string {
	return "N"
}

func (LTDD) Name() string {
	return "旅途点点（田螺科技）"
}

func (LTDD) URL() (pattern string, mockValue interface{}) {
	return "https://wechat.lvtudiandian.com/index.php/Home/SweepCode/index?locomotiveId=%s&carriage=%d&seatRow=%d&seatNo=%v", "D/F"
}

func (LTDD) AlwaysOn() bool {
	return true
}

func (a LTDD) RoundTrip(req *http.Request) (*http.Response, error) {
	common.SetCookies(req, []*http.Cookie{{
		Name:  "OpenId",
		Value: SessionID(a),
	}})
	return AdapterConf(a).Request.RoundTrip(req)
}

func (a LTDD) Info(serial string) (info JSONObject, err error) {
	const (
		landingPage  = "https://wechat.lvtudiandian.com/index.php/QrSweepCode/index?locomotiveId=%s&carriage=6&seatRow=6&seatNo=D%%2FF"
		orderingPage = "https://wechat.lvtudiandian.com/index.php/Home/SweepCode/index.html?is_redirect=1"
	)

	url := fmt.Sprintf(landingPage, serial)
	var resp *http.Response
	if resp, err = common.HTTPClient(a).Get(url); err != nil {
		return
	}
	defer resp.Body.Close()

	var bytes []byte
	if bytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if strings.HasPrefix(string(bytes), "<script>alert") {
		return
	}

	if resp, err = common.HTTPClient(a).Get(orderingPage); err != nil {
		return
	}
	defer resp.Body.Close()

	if bytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	if match := jsonUnitNoRegExp.FindSubmatch(bytes); match != nil {
		json.Unmarshal(match[1], &info)
	} else if match := htmlUnitNoRegExp.FindSubmatch(bytes); match != nil {
		info = JSONObject{"locomotive_code": string(match[1])}
		if match = htmlTrainNoRegExp.FindSubmatch(bytes); match != nil {
			info["partner_name"] = string(match[1])
		}
	}
	if match := jsonCompanyRegExp.FindSubmatch(bytes); match != nil {
		json.Unmarshal(match[1], &info)
	}
	return
}

func (LTDD) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	trains = []TrainSchedule{{
		TrainNo: info["partner_name"].(string),
	}}
	return
}

func (LTDD) UnitNo(_ string, info JSONObject) (unitNo string, err error) {
	defer common.Catch(&err)
	unitNo = common.NormalizeUnitNo(info["locomotive_code"].(string))
	if strings.HasPrefix(unitNo, "380") {
		unitNo = "CRH" + unitNo
	}
	return
}

func (LTDD) Operator(_ string, info JSONObject) (bureauCode string, err error) {
	defer common.Catch(&err)
	bureauCode = jsonCompanyIDMap[info["company_id"].(string)]
	return
}
