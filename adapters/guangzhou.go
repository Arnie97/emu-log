package adapters

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/arnie97/emu-log/common"
)

type Guangzhou struct{}

func init() {
	Register(Guangzhou{})
}

func (Guangzhou) Code() string {
	return "Q"
}

func (Guangzhou) Name() string {
	return "中国铁路广州局集团有限公司"
}

func (Guangzhou) URL() string {
	return "http://sj.yishizongheng.com/?code=%s"
}

func (Guangzhou) BruteForce(serials chan<- string) {
	for x := 1; x < 90; x++ {
		serials <- fmt.Sprintf("%03d", x)
	}
	for x := 220; x < 600; x++ {
		serials <- fmt.Sprintf("%03d", x)
	}
}

func (b Guangzhou) Info(serial string) (info jsonObject, err error) {
	const api = "https://sj-api.yishizongheng.com/shejian/api/train/getByQrcode?qrcode=%s"
	url := fmt.Sprintf(api, strings.TrimLeft(serial, "0"))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header.Set("authorization", common.Conf(b.Code()))
	resp, err := common.HTTPClient().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status int    `json:"code"`
		Msg    string `json:"message"`
		Data   jsonObject
	}
	err = parseResult(resp, &result)
	info = result.Data
	if info != nil {
		err = nil
		info["serial"] = serial
	}
	return
}

func (b Guangzhou) TrainNo(info jsonObject) (trainNo, date string, err error) {
	defer common.Catch(&err)
	trainNo = info["train"].(string)
	return
}

func (b Guangzhou) VehicleNo(info jsonObject) (vehicleNo string, err error) {
	defer common.Catch(&err)
	vehicleNo = fmt.Sprintf(
		"CR%s-%.0f+%s",
		info["carriageNum"].(string), info["id"], info["serial"].(string),
	)
	return
}
