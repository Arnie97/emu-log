package adapters

import (
	"fmt"
	"net/http"
	"strings"
	"time"

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

func (Zhengzhou) AlwaysOn() bool {
	return true
}

func (b Zhengzhou) RoundTrip(req *http.Request) (*http.Response, error) {
	time.Sleep(common.RequestInterval)
	req.Header.Set("user-agent", common.UserAgentJDPay)
	if len(req.Cookies()) == 0 {
		req.Header.Set("cookie", common.Conf(b.Code()))
	}
	return http.DefaultTransport.RoundTrip(req)
}

func (b Zhengzhou) Info(serial string) (info jsonObject, err error) {
	const api = "https://p.12306.cn/tservice/mealAction/qrcodeDecode"
	var (
		cookie *http.Cookie
		req    *http.Request
		resp   *http.Response
	)
	if cookie, err = b.OAuth(serial); err != nil {
		return
	}
	if req, err = http.NewRequest(http.MethodPost, api, nil); err != nil {
		return
	}
	req.Header.Set("cookie", cookie.Name+"="+cookie.Value)
	if resp, err = common.HTTPClient(b).Do(req); err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status bool   `json:"status"`
		Msg    string `json:"errorMsg"`
		Data   struct {
			TrainQrcodeInfo jsonObject
		}
	}
	err = parseResult(resp, &result)
	info = result.Data.TrainQrcodeInfo
	return
}

func (b Zhengzhou) OAuth(serial string) (cookie *http.Cookie, err error) {
	var resp *http.Response
	if resp, err = common.HTTPClient(b).Get(BuildURL(b, serial)); err != nil {
		return
	}

	authURL := resp.Request.URL
	authURL.Path = authURL.Path + "_fbs"
	if resp, err = common.HTTPClient(b).Get(authURL.String()); err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status int    `json:"errcode"`
		Msg    string `json:"errmsg"`
		URL    string `json:"return_url"`
	}
	if err = parseResult(resp, &result); err != nil {
		return
	}
	if resp, err = common.HTTPClient(b).Get(result.URL); err != nil {
		return
	}
	defer resp.Body.Close()
	for _, each := range resp.Cookies() {
		if each.Name == "JSESSIONID" {
			cookie = each
			return
		}
	}
	err = fmt.Errorf("failed to acquire user session ID with OAuth")
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
	if strings.HasPrefix(vehicleNo, "CHR") {
		vehicleNo = strings.Replace(vehicleNo, "CHR", "CRH", 1)
	}
	return
}
