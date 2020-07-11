// Package adapters defines API adapters for each supported railway bureau.
package adapters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

type (
	Bureau interface {
		Code() string
		Name() string
		URL() string
		BruteForce(chan<- string)
		Info(serial string) (info jsonObject, err error)
		TrainNo(info jsonObject) (trainNo, date string, err error)
		VehicleNo(info jsonObject) (vehicleNo string, err error)
	}
	jsonObject map[string]interface{}
)

var (
	Bureaus = make(map[string]Bureau)
)

func Register(b Bureau) {
	if Bureaus[b.Code()] != nil {
		common.Must(fmt.Errorf("[%s] bureau was redeclared in the adapters package", b.Code()))
	}
	Bureaus[b.Code()] = b
	log.Debug().Msgf("[%s] bureau adapter registered: %s", b.Code(), b.Name())
}

func MustGetBureauByCode(bureauCode string) (b Bureau) {
	if b = Bureaus[bureauCode]; b == nil {
		log.Fatal().Msgf("[%s] unknown bureau adapter", bureauCode)
	}
	return
}

func BuildURL(b Bureau, serial string) (url string) {
	return fmt.Sprintf(b.URL(), serial)
}

func ParseURL(url string) (b Bureau, serial string) {
	for _, b = range Bureaus {
		urlPattern := fmt.Sprintf(regexp.QuoteMeta(b.URL()), `(\w+)`)
		urlRegExp := regexp.MustCompile(urlPattern)
		if match := urlRegExp.FindStringSubmatch(url); match != nil {
			return b, match[1]
		}
	}
	return nil, ""
}

func parseResult(resp *http.Response, resultPtr interface{}) (err error) {
	err = json.NewDecoder(resp.Body).Decode(resultPtr)
	if err != nil {
		return
	}

	var (
		ok     bool
		status = common.GetField(resultPtr, "Status")
		msg    = common.GetField(resultPtr, "Msg")
	)
	switch status.(type) {
	case string:
		ok = status == "ok"
	case int:
		ok = status == 200 || status == 0
	default:
		ok = false
	}
	if !ok {
		err = fmt.Errorf("api error %v: %s", status, msg)
	}
	return
}
