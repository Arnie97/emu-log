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

func (Beijing) BruteForce(qrCodes chan<- string) {
	for y := 1; y <= 4; y++ {
		for x := 11000; x < 16000; x += 500 {
			qrCodes <- fmt.Sprintf("%d%07d", y, x)
		}
	}
	for x := 342000; x < 640000; x += 500 {
		qrCodes <- fmt.Sprintf("5%07d", x)
	}
	for x := 1000; x < 500000; x += 500 {
		qrCodes <- fmt.Sprintf("6%07d", x)
	}
	for y := 7; y <= 9; y++ {
		for x := 11000; x < 16000; x += 500 {
			qrCodes <- fmt.Sprintf("%d%07d", y, x)
		}
	}
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

func (b Beijing) TrainNo(qrCode string) (trainNo, date string, err error) {
	var info jsonObject
	info, err = b.Info(qrCode)
	if err == nil {
		defer common.Catch(&err)
		trainNo = info["TrainnoId"].(string)
		date = info["TrainnoDate"].(string)
	}
	return
}

func (b Beijing) VehicleNo(qrCode string) (vehicleNo string, err error) {
	var info jsonObject
	info, err = b.Info(qrCode)
	if err == nil {
		defer common.Catch(&err)
		vehicleNo = common.NormalizeVehicleNo(info["TrainId"].(string))
	}
	return
}