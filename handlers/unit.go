package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/arnie97/emu-log/models"
	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
)

type urlWrapper struct {
	URL *string `json:"url,omitempty"`
}

// singleUnitNoHandler takes an exact unit number,
// and returns the 30 most recent log items for the unit.
func singleUnitNoHandler(w http.ResponseWriter, r *http.Request) {
	unitNo := common.NormalizeUnitNo(chi.URLParam(r, "unitNo"))
	results := models.ListTrainsForSingleUnitNo("%" + unitNo)
	jsonResponse(results, w)
}

// multiUnitNoHandler takes multiple exact unit numbers,
// and returns the most recent log item for each of them.
func multiUnitNoHandler(w http.ResponseWriter, r *http.Request) {
	unitNoList := strings.Split(chi.URLParam(r, "unitNo"), ",")
	results := models.ListLatestTrainForMultiUnits(unitNoList)
	jsonResponse(results, w)
}

// fuzzyUnitNoHandler takes an incomplete part of the unit number,
// and returns the most recent occurance for the first 30 units
// in lexicographical order that matches the given fuzzy pattern.
func fuzzyUnitNoHandler(w http.ResponseWriter, r *http.Request) {
	unitNo := common.NormalizeUnitNo(chi.URLParam(r, "unitNo"))
	results := models.ListLatestTrainForMatchedUnits(unitNo)
	jsonResponse(results, w)
}

// unitBuildURLHandler takes an exact unit number, and rebuild
// the URL encoded in one of the QR code stickers attached to the unit.
func unitBuildURLHandler(w http.ResponseWriter, r *http.Request) {
	unitNo := common.NormalizeUnitNo(chi.URLParam(r, "unitNo"))
	serialModels := models.ListSerialsForSingleUnit(unitNo)
	var results urlWrapper
	for _, s := range serialModels {
		if a := adapters.Adapters[s.Adapter]; a != nil {
			url := adapters.BuildURL(a, s.SerialNo)
			results.URL = &url
			break
		}
	}
	jsonResponse(results, w)
}

// unitParseURLHandler tries to parse the URL
// encoded in one of the QR code stickers attached to the unit.
func unitParseURLHandler(w http.ResponseWriter, r *http.Request) {
	// case 1: invalid URL
	var (
		input   urlWrapper
		results = make([]models.LogModel, 0)
	)
	defer func() {
		jsonResponse(results, w)
	}()
	if json.NewDecoder(r.Body).Decode(&input); input.URL == nil {
		w.WriteHeader(400)
		return
	}

	// case 2: unrecognized URL
	a, serial := adapters.ParseURL(*input.URL)
	unitNo := common.NormalizeUnitNo(chi.URLParam(r, "unitNo"))
	if a == nil {
		serialModel := models.SerialModel{
			Adapter:  "?",
			Operator: "?",
			UnitNo:   unitNo,
			SerialNo: *input.URL,
		}
		serialModel.Add()
		log.Debug().Msgf("[%s] %v", serialModel.Adapter, serialModel)
		return
	}

	// case 3: unit currently offline or unit number does not match
	serialModel := models.SerialModel{
		Adapter:  a.Code(),
		Operator: "?",
		SerialNo: serial,
	}
	info, err := a.Info(serial)
	if err == nil {
		serialModel.UnitNo, err = a.UnitNo(serial, info)
	}
	serialModel.Operator, err = adapters.Operator(a, serial, info)
	if !common.ApproxEqualUnitNo(unitNo, serialModel.UnitNo) {
		serialModel.UnitNo = "-" + unitNo + "@" + serialModel.UnitNo
		serialModel.Add()
		log.Debug().Msgf("[%s] %v", serialModel.Adapter, serialModel)
		return
	}

	// case 4: unit number matches user input
	serialModel.Add()
	serialModel.AddTrainOperationLogs(info)

	results = models.ListTrainsForSingleUnitNo(serialModel.UnitNo)
	return
}

// unitParseURLMapHandler tries to parse the URL in the QR code stickers
// attached to the unit and returns a map as result.
func unitParseURLMapHandler(w http.ResponseWriter, r *http.Request) {
	// case 1: invalid URL
	var (
		input   urlWrapper
		mapResp interface{}
	)
	defer func() {
		jsonResponse(mapResp, w)
	}()
	if json.NewDecoder(r.Body).Decode(&input); input.URL == nil {
		w.WriteHeader(400)
		return
	}

	// case 2: unrecognized URL
	a, serial := adapters.ParseURL(*input.URL)
	if a == nil {
		return
	}

	// case 3: unit currently offline
	serialModel := models.SerialModel{
		Adapter:  a.Code(),
		SerialNo: serial,
	}

	info, err := a.Info(serial)
	if err == nil {
		serialModel.UnitNo, err = a.UnitNo(serial, info)
	}
	serialModel.Operator, err = adapters.Operator(a, serial, info)
	serialModel.Add()
	serialModel.AddTrainOperationLogs(info)

	resp := struct {
		Operator string `json:"operator,omitempty"`
		UnitNo   string `json:"emu_no,omitempty"`
		SerialNo string `json:"serial_no"`
		Adapter  struct {
			Code string `json:"code"`
			Name string `json:"name"`
		} `json:"adapter"`
		Logs []models.LogModel `json:"logs,omitempty"`
	}{
		Operator: serialModel.Operator,
		UnitNo:   serialModel.UnitNo,
		SerialNo: serial,
		Logs:     models.ListTrainsForSingleUnitNo(serialModel.UnitNo),
	}
	resp.Adapter.Code = a.Code()
	resp.Adapter.Name = a.Name()
	mapResp = resp

	log.Debug().Msgf("[%s] %+v", serialModel.Adapter, mapResp)
	return
}
