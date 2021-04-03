package adapters

import (
	"net/http"
	"strings"
	"time"

	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

type Zhengzhou struct {
	cookies []*http.Cookie
}

func init() {
	Register(&Zhengzhou{})
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

func (b *Zhengzhou) RoundTrip(req *http.Request) (*http.Response, error) {
	time.Sleep(common.RequestInterval)
	common.SetUserAgent(req, common.UserAgentJDPay)
	common.SetCookies(req, b.cookies)
	resp, err := http.DefaultTransport.RoundTrip(req)

	// stop further redirects and collect crucial cookies
	if err == nil && resp != nil && resp.StatusCode == http.StatusFound {
		switch req.URL.Path {
		case "/cgi-bin/app/appjmp":
			log.Info().Msgf("%v %+v", resp.Status, resp.Cookies())
			fallthrough
		case "/tservice/catering/jd":
			b.cookies = resp.Cookies()
			resp.StatusCode = http.StatusOK
		}
	}
	return resp, err
}

func (b *Zhengzhou) Info(serial string) (info JSONObject, err error) {
	var resp *http.Response
	if resp, err = b.OAuth(serial); err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status bool   `json:"status"`
		Msg    string `json:"errorMsg"`
		Data   struct {
			TrainQrcodeInfo JSONObject
		}
	}
	err = parseResult(resp, &result)
	info = result.Data.TrainQrcodeInfo
	return
}

// RefreshToken applies for a new access token from JD pay
// if the cached access token has already been expired.
func (b *Zhengzhou) RefreshToken() (resp *http.Response, err error) {
	if len(b.cookies) > 0 {
		return
	}

	// the access tokens will be saved by the custom round tripper
	const api = "https://ms.jr.jd.com/jrmserver/base/user/getNewTokenJumpUrl"
	return common.HTTPClient(b).Get(api + common.Conf(b.Code()))
}

// OAuth obtains a new authorization code from JD pay,
// and start a new session on 12306 servers with it.
// This is required for each run, since each session is bound
// to an immutable vehicle serial number.
func (b *Zhengzhou) OAuth(serial string) (resp *http.Response, err error) {
	if resp, err = b.RefreshToken(); err != nil {
		return
	}
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
		// refresh the expiry access token and try again
		if result.Status == 256 {
			b.cookies = nil
			return b.OAuth(serial)
		}
		return
	}

	// fork into a new session
	session := *b
	// the session ID will be saved by the custom round tripper
	if resp, err = common.HTTPClient(&session).Get(result.URL); err != nil {
		return
	}

	const api = "https://p.12306.cn/tservice/mealAction/qrcodeDecode"
	return common.HTTPClient(&session).PostForm(api, nil)
}

func (Zhengzhou) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	shortDate := info["startDay"].(string)
	trains = []TrainSchedule{{
		TrainNo: info["trainCode"].(string),
		Date:    shortDate[:4] + "-" + shortDate[4:6] + "-" + shortDate[6:8],
	}}
	return
}

func (Zhengzhou) VehicleNo(info JSONObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	vehicleNo = common.NormalizeVehicleNo(info["carCode"].(string))
	if strings.HasPrefix(vehicleNo, "CHR") {
		vehicleNo = strings.Replace(vehicleNo, "CHR", "CRH", 1)
	}
	return
}
