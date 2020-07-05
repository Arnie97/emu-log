package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

func (Guangzhou) BruteForce(serials chan<- string) {
	for x := 1; x < 90; x++ {
		serials <- fmt.Sprintf("%03d", x)
	}
}

func (Guangzhou) Info(serial string) (info jsonObject, err error) {
	const api = "https://v3i.minicart.cn/shopping/v3/getTrainnum"
	const contentType = "application/json"
	values := jsonObject{
		"qr_code": strings.TrimLeft(serial, "0"),
		"mpid":    9,
		"mp_id":   9,
		"mid":     9,
		"token":   "2107e4f9dc309b5f8a5b05b9b7483cc0",
	}
	jsonStr, err := json.Marshal(values)
	if err != nil {
		return
	}
	resp, err := common.HTTPClient().Post(api, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"error_code"`
		Msg    string
		Data   jsonObject
	}
	err = parseResult(resp, &result)
	info = result.Data
	return
}

func (b Guangzhou) TrainNo(serial string) (trainNo, date string, err error) {
	var info jsonObject
	info, err = b.Info(serial)
	if err == nil {
		defer common.Catch(&err)
		trainNo = info["train"].(string)
		date = time.Now().Format("2006-01-02")
	}
	return
}

func (b Guangzhou) VehicleNo(serial string) (vehicleNo string, err error) {
	var info jsonObject
	info, err = b.Info(serial)
	if err == nil {
		defer common.Catch(&err)
		vehicleNo = fmt.Sprintf("CR%s+%s", info["carriage_num"], serial)
	}
	return
}
