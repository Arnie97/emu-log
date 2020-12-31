package adapters

import (
	"net/http"
	"net/url"
	"strings"

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

func (b Zhengzhou) Info(serial string) (info jsonObject, err error) {
	const api = "https://p.12306.cn/tservice/mealAction/qrcodeDecode"
	if err = b.OAuth(serial); err != nil {
		return
	}

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
		Status bool   `json:"status"`
		Msg    string `json:"errMsg"`
		Data   struct {
			TrainQrcodeInfo jsonObject
		}
	}
	err = parseResult(resp, &result)
	info = result.Data.TrainQrcodeInfo
	return
}

func (b Zhengzhou) OAuth(serial string) (err error) {
	req, err := http.NewRequest("GET", b.AuthURL(serial), nil)
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
		Status int    `json:"errcode"`
		Msg    string `json:"errmsg"`
		URL    string `json:"return_url"`
	}
	if err = parseResult(resp, &result); err != nil {
		return
	}
	if req, err = http.NewRequest("GET", result.URL, nil); err != nil {
		return
	}
	req.Header.Set("Cookie", common.Conf(b.Code()))
	if resp, err = common.HTTPClient().Do(req); err != nil {
		return
	}
	defer resp.Body.Close()
	return
}

func (b Zhengzhou) AuthURL(serial string) (authURL string) {
	authURL = strings.Replace(BuildURL(b, serial), "/init", "/jd", 1)
	authURL = url.QueryEscape(url.QueryEscape(authURL))
	authURL = "https://mobile.12306.cn/weixin/jd/auth?redirect=" + authURL
	authURL = "https://jauth.jd.com/entrance_fbs?" + url.Values{
		"response_type": {"code"},
		"appid":         {"jd8c6431caca1f6602"},
		"scope":         {"scope.mobile,scope.userInfo"},
		"redirect_uri":  {authURL},
		"cancel_uri":    {""},
		"act_type":      {"2"},
		"state":         {"12306weixin"},
		"show_title":    {"1"},
	}.Encode()
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
