package adapters

import (
	"net/http"

	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

type JD12306 struct {
	cookies []*http.Cookie
}

func init() {
	Register(&JD12306{})
}

func (JD12306) Code() string {
	return "F"
}

func (JD12306) Name() string {
	return "京东金融"
}

func (JD12306) URL() (pattern string, mockValue interface{}) {
	return "https://p.12306.cn/tservice/catering/init?c=%s&w=%v", "t"
}

func (JD12306) AlwaysOn() bool {
	return true
}

func (a *JD12306) RoundTrip(req *http.Request) (*http.Response, error) {
	common.SetCookies(req, a.cookies)
	resp, err := AdapterConf(a).Request.RoundTrip(req)

	// stop further redirects and collect crucial cookies
	if err == nil && resp != nil && resp.StatusCode == http.StatusFound {
		switch req.URL.Path {
		case "/cgi-bin/app/appjmp":
			log.Info().Msgf("%v %+v", resp.Status, resp.Cookies())
			fallthrough
		case "/tservice/catering/jd":
			a.cookies = resp.Cookies()
			resp.StatusCode = http.StatusOK
		}
	}
	return resp, err
}

func (a *JD12306) Info(serial string) (info JSONObject, err error) {
	var resp *http.Response
	if resp, err = a.OAuth(serial); err != nil {
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
func (a *JD12306) RefreshToken() (resp *http.Response, err error) {
	if len(a.cookies) > 0 {
		return
	}

	// the access tokens will be saved by the custom round tripper
	const api = "https://ms.jr.jd.com/jrmserver/base/user/getNewTokenJumpUrl"
	return common.HTTPClient(a).Get(api + SessionID(a))
}

// OAuth obtains a new authorization code from JD pay,
// and start a new session on 12306 servers with it.
// This is required for each run, since each session is bound
// to an immutable unit serial number.
func (a *JD12306) OAuth(serial string) (resp *http.Response, err error) {
	if resp, err = a.RefreshToken(); err != nil {
		return
	}
	if resp, err = common.HTTPClient(a).Get(BuildURL(a, serial)); err != nil {
		return
	}

	authURL := resp.Request.URL
	authURL.Path = authURL.Path + "_fbs"
	if resp, err = common.HTTPClient(a).Get(authURL.String()); err != nil {
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
			a.cookies = nil
			return a.OAuth(serial)
		}
		return
	}

	// fork into a new session
	session := *a
	// the session ID will be saved by the custom round tripper
	if resp, err = common.HTTPClient(&session).Get(result.URL); err != nil {
		return
	}

	const api = "https://p.12306.cn/tservice/mealAction/qrcodeDecode"
	return common.HTTPClient(&session).PostForm(api, nil)
}

func (JD12306) TrainNo(info JSONObject) (trains []TrainSchedule, err error) {
	defer common.Catch(&err)
	shortDate := info["startDay"].(string)
	trains = []TrainSchedule{{
		TrainNo: info["trainCode"].(string),
		Date:    shortDate[:4] + "-" + shortDate[4:6] + "-" + shortDate[6:8],
	}}
	return
}

func (JD12306) UnitNo(_ string, info JSONObject) (unitNo string, err error) {
	defer common.Catch(&err)
	unitNo = common.NormalizeUnitNo(info["carCode"].(string))
	return
}

func (JD12306) Operator(serialNo string, _ JSONObject) (bureauCode string, err error) {
	defer common.Catch(&err)
	bureauCode = serialNo[:1]
	return
}
