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

		// URL returns the URL pattern contained in the QR codes,
		// with the serial number replaced by the placeholder "%s".
		URL() string

		// BruteForce takes a channel, and sends all possibly valid
		// serial numbers into the channel in lexicographical order.
		BruteForce(serialNo chan<- string)

		// AlwaysOn means the bureau adapter still returns some basic
		// information even if meal ordering service is currently not
		// available. Otherwise, unallocated serial numbers cannot
		// be differentiated from serials assigned to offline vehicles,
		// and the unknown serials have to be visited in each scan.
		AlwaysOn() bool

		Info(serialNo string) (info jsonObject, err error)
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
	switch status := status.(type) {
	case int:
		ok = status == 0 || status == 200 || status == 257
	case bool:
		ok = status
	default:
		ok = false
	}
	if !ok {
		err = fmt.Errorf("api error %v: %s", status, msg)
	}
	return
}
