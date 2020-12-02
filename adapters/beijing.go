package adapters

import (
	"crypto/md5"
	"fmt"
	"net/url"

	"github.com/arnie97/emu-log/common"
)

type Beijing struct{}

func init() {
	Register(Beijing{})
}

func (Beijing) Code() string {
	return "P"
}

func (Beijing) Name() string {
	return "中国铁路北京局集团有限公司"
}

func (Beijing) URL() string {
	return "https://aymaoto.jtlf.cn/page/oto/index?QR=%s"
}

var (
	turn  int
	shift = []int{0, 1000, 500, 1500}
)

func (Beijing) BruteForce(qrCodes chan<- string) {
	turn = (turn + 1) % 4
	for x := shift[turn]; x < 990000; x += 2000 {
		qrCodes <- fmt.Sprintf("5%07d", x)
	}
	for x := shift[turn]; x < 700000; x += 2000 {
		qrCodes <- fmt.Sprintf("6%07d", x)
	}
}

func (Beijing) AlwaysOn() bool {
	return false
}

func (Beijing) Info(qrCode string) (info jsonObject, err error) {
	const api = "https://aymaoto.jtlf.cn/webapi/otoshopping/ewh_getqrcodetrainnoinfo"
	const key = "qrcode=%s&key=ltRsjkiM8IRbC80Ni1jzU5jiO6pJvbKd"
	sign := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf(key, qrCode))))
	form := url.Values{"qrCode": {qrCode}, "sign": {sign}}
	resp, err := common.HTTPClient().PostForm(api, form)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status int `json:"state"`
		Msg    string
		Data   struct {
			TrainInfo jsonObject
			URLStr    string
		}
	}
	err = parseResult(resp, &result)
	info = result.Data.TrainInfo
	return
}

func (b Beijing) TrainNo(info jsonObject) (trainNo, date string, err error) {
	defer common.Catch(&err)
	trainNo = info["TrainnoId"].(string)
	date = info["TrainnoDate"].(string)
	return
}

func (b Beijing) VehicleNo(info jsonObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	vehicleNo = common.NormalizeVehicleNo(info["TrainId"].(string))
	return
}
