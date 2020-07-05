// Package adapters defines API adapters for each supported railway bureau.
package adapters

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

type (
	Bureau interface {
		Code() string
		Name() string
		BruteForce(chan<- string)
		Info(qrCode string) (info jsonObject, err error)
		TrainNo(qrCode string) (trainNo, date string, err error)
		VehicleNo(qrCode string) (vehicleNo string, err error)
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
		ok = status.(string) == "ok"
	case int:
		ok = status.(int) == 200
	default:
		ok = false
	}
	if !ok {
		err = fmt.Errorf("api error %v: %s", status, msg)
	}
	return
}
