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

func (b Wuhan) Info(serial string) (info jsonObject, err error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(b.URL(), serial), nil)
	if err != nil {
		return
	}
	req.Header.Set("Cookie", common.Conf(b.Code()))
	resp, err := common.HTTPClient().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if match := jsonVehicleRegExp.FindSubmatch(bytes); match != nil {
		json.Unmarshal(match[1], &info)
	} else if match := htmlVehicleRegExp.FindSubmatch(bytes); match != nil {
		info = jsonObject{"locomotive_code": string(match[1])}
	}
	if match := jsonCompanyRegExp.FindSubmatch(bytes); match != nil {
		json.Unmarshal(match[1], &info)
	}
	return
}

func (b Wuhan) TrainNo(info jsonObject) (trainNo, date string, err error) {
	defer common.Catch(&err)
	trainNo = info["partner_name"].(string)
	return
}

func (b Wuhan) VehicleNo(info jsonObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	vehicleNo = common.NormalizeVehicleNo(info["locomotive_code"].(string))
	if strings.HasPrefix(vehicleNo, "380") {
		vehicleNo = "CRH" + vehicleNo
	} else if strings.HasPrefix(vehicleNo, "CRH400") {
		vehicleNo = strings.Replace(vehicleNo, "CRH", "CR", 1)
	}
	return
}
