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
	htmlVehicleRegExp = regexp.MustCompile(`<p class="schedule">.+:(\w+)-6-6D/F</p>`)
	htmlTrainNoRegExp = regexp.MustCompile(`<p class="train-title">(\w+)随车购</p>`)
	jsonVehicleRegExp = regexp.MustCompile(`var locomotive_info = (\{.+\});`)
	jsonCompanyRegExp = regexp.MustCompile(`var company_info = (\{.+\});`)
)

type Wuhan struct{}

func init() {
	Register(Wuhan{})
}

func (Wuhan) Code() string {
	return "N"
}

func (Wuhan) Name() string {
	return "中国铁路武汉局集团有限公司"
}

func (Wuhan) URL() string {
	return "https://wechat.lvtudiandian.com/index.php/Home/SweepCode/index?locomotiveId=%s&carriage=6&seatRow=6&seatNo=D/F"
}

func (Wuhan) BruteForce(serials chan<- string) {
	for x := 1; x < 500; x++ {
		serials <- fmt.Sprintf("%03d", x)
	}
}

func (Wuhan) AlwaysOn() bool {
	return true
}

func (b Wuhan) RoundTrip(req *http.Request) (*http.Response, error) {
	common.SetCookies(req, []*http.Cookie{{
		Name:  "OpenId",
		Value: common.Conf(b.Code()),
	}})
	return common.IntervalTransport{}.RoundTrip(req)
}

func (b Wuhan) Info(serial string) (info JSONObject, err error) {
	const (
		landingPage  = "https://wechat.lvtudiandian.com/index.php/QrSweepCode/index?locomotiveId=%s&carriage=6&seatRow=6&seatNo=D%%2FF"
		orderingPage = "https://wechat.lvtudiandian.com/index.php/Home/SweepCode/index.html?is_redirect=1"
	)

	url := fmt.Sprintf(landingPage, serial)
	var resp *http.Response
	if resp, err = common.HTTPClient(b).Get(url); err != nil {
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

	if resp, err = common.HTTPClient(b).Get(orderingPage); err != nil {
		return
	}
	defer resp.Body.Close()

	if bytes, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	if match := jsonVehicleRegExp.FindSubmatch(bytes); match != nil {
		json.Unmarshal(match[1], &info)
	} else if match := htmlVehicleRegExp.FindSubmatch(bytes); match != nil {
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

func (Wuhan) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	trains = []TrainSchedule{{
		TrainNo: info["partner_name"].(string),
	}}
	return
}

func (Wuhan) VehicleNo(info JSONObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	vehicleNo = common.NormalizeVehicleNo(info["locomotive_code"].(string))
	if strings.HasPrefix(vehicleNo, "380") {
		vehicleNo = "CRH" + vehicleNo
	} else if strings.HasPrefix(vehicleNo, "CRH400") {
		vehicleNo = strings.Replace(vehicleNo, "CRH", "CR", 1)
	}
	return
}
