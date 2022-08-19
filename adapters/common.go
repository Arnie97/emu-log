// Package adapters defines API adapters for each supported site.
package adapters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

type (
	Adapter interface {
		Code() string
		Name() string

		// URL returns the URL pattern contained in the QR codes,
		// with the serial number replaced by the placeholder "%s".
		URL() (pattern string, mockValue interface{})

		// AlwaysOn means the site adapter still returns some basic
		// information even if meal ordering service is currently not
		// available. Otherwise, unallocated serial numbers cannot
		// be differentiated from serials assigned to offline units,
		// and the unknown serials have to be visited in each scan.
		AlwaysOn() bool

		Info(serialNo string) (info JSONObject, err error)
		TrainNo(info JSONObject) ([]TrainSchedule, error)
		UnitNo(serialNo string, info JSONObject) (unitNo string, err error)
	}
	UnionAdapter interface {
		Adapter
		Operator(serialNo string, info JSONObject) (bureauCode string, err error)
	}
	JSONObject    map[string]interface{}
	TrainSchedule struct{ TrainNo, Date string }
)

var (
	printfMap = strings.NewReplacer(
		"%s", "%[1]v",
		"%d", "%[2]v",
		"%v", "%[3]v",
		"%02d", "%02[2]v",
	)
	printfRegExps = []interface{}{
		`([-\w\s,]+)`, // %s
		`(?:\d+)`,     // %d
		`(?:.*?)`,     // %v
	}
	Adapters = make(map[string]Adapter)
)

func Register(a Adapter) {
	if Adapters[a.Code()] != nil {
		common.Must(fmt.Errorf("[%s] adapter was redeclared in the adapters package", a.Code()))
	}
	Adapters[a.Code()] = a
}

func MustGetAdapterByCode(adapterCode string) (a Adapter) {
	if a = Adapters[adapterCode]; a == nil {
		log.Fatal().Msgf("[%s] unknown adapter", adapterCode)
	}
	return
}

func AdapterConf(a Adapter) (merged common.AdapterConf) {
	global := common.Conf()
	merged = global.Adapters[a.Code()]
	adapterRequest := merged.Request
	merged.Request = new(common.RequestConf)
	common.Must(common.StructDecode(global.Request, merged.Request))
	common.Must(common.StructDecode(adapterRequest, merged.Request))
	return
}

func SessionID(a Adapter) string {
	return AdapterConf(a).Request.SessionID
}

func Operator(a Adapter, serial string, info JSONObject) (bureauCode string, err error) {
	if a, ok := a.(UnionAdapter); ok {
		return a.Operator(serial, info)
	}
	return a.Code(), nil
}

func BuildURL(a Adapter, serial string) (url string) {
	urlPattern, urlMockValue := a.URL()
	return fmt.Sprintf(printfMap.Replace(urlPattern), serial, 6, urlMockValue)
}

func ParseURL(url string) (a Adapter, serial string) {
	for _, a = range Adapters {
		urlPattern, _ := a.URL()
		urlPattern = printfMap.Replace(regexp.QuoteMeta(urlPattern))
		urlRegExp := regexp.MustCompile(fmt.Sprintf(urlPattern, printfRegExps...))
		if match := urlRegExp.FindStringSubmatch(url); match != nil {
			return a, match[1]
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
