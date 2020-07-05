package adapters

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/arnie97/emu-log/common"
)

var (
	extractRegExp = regexp.MustCompile(`<p class="schedule">.+?:(\w+)-*</p>`)
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

func (Wuhan) BruteForce(serials chan<- string) {
	for x := 1; x < 500; x++ {
		serials <- fmt.Sprintf("%03d", x)
	}
}

func (Wuhan) Info(serial string) (info jsonObject, err error) {
	const api = "https://wechat.lvtudiandian.com/index.php/Home/SweepCode/index?locomotiveId=%s"
	resp, err := common.HTTPClient().Get(fmt.Sprintf(api, serial))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	extract := extractRegExp.FindSubmatch(bytes)
	if len(extract) < 2 {
		return
	}

	key := "vehicleNo"
	if len(common.NormalizeTrainNo(string(extract[1]))) != 0 {
		key = "trainNo"
	}
	info = jsonObject{key: string(extract[1])}
	return
}

func (b Wuhan) TrainNo(serial string) (trainNo, date string, err error) {
	var info jsonObject
	info, err = b.Info(serial)
	if err == nil {
		defer common.Catch(&err)
		trainNo = info["trainNo"].(string)
	}
	return
}

func (b Wuhan) VehicleNo(serial string) (vehicleNo string, err error) {
	var info jsonObject
	info, err = b.Info(serial)
	if err == nil {
		defer common.Catch(&err)
		vehicleNo = info["vehicleNo"].(string)
	}
	return
}
