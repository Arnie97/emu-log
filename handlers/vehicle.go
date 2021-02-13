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

// singleVehicleNoHandler takes an exact vehicle number,
// and returns the 30 most recent log items for the vehicle.
func singleVehicleNoHandler(w http.ResponseWriter, r *http.Request) {
	vehicleNo := common.NormalizeVehicleNo(chi.URLParam(r, "vehicleNo"))
	results := models.ListTrainsForSingleVehicle("%" + vehicleNo)
	jsonResponse(results, w)
}

// multiVehicleNoHandler takes multiple exact vehicles numbers,
// and returns the most recent log item for each of them.
func multiVehicleNoHandler(w http.ResponseWriter, r *http.Request) {
	vehicleNoList := strings.Split(chi.URLParam(r, "vehicleNo"), ",")
	results := models.ListLatestTrainForMultiVehicles(vehicleNoList)
	jsonResponse(results, w)
}

// fuzzyVehicleNoHandler takes an incomplete part of the vehicle number,
// and returns the most recent occurance for the first 30 vehicles
// in lexicographical order that matches the given fuzzy pattern.
func fuzzyVehicleNoHandler(w http.ResponseWriter, r *http.Request) {
	vehicleNo := common.NormalizeVehicleNo(chi.URLParam(r, "vehicleNo"))
	results := models.ListLatestTrainForMatchedVehicles(vehicleNo)
	jsonResponse(results, w)
}

// vehicleBuildURLHandler takes an exact vehicle number, and rebuild
// the URL encoded in one of the QR code stickers attached to the vehicle.
func vehicleBuildURLHandler(w http.ResponseWriter, r *http.Request) {
	vehicleNo := common.NormalizeVehicleNo(chi.URLParam(r, "vehicleNo"))
	serialModels := models.ListSerialsForSingleVehicle(vehicleNo)
	var results urlWrapper
	for _, s := range serialModels {
		if b := adapters.Bureaus[s.BureauCode]; b != nil {
			url := adapters.BuildURL(b, s.SerialNo)
			results.URL = &url
		}
	}
	jsonResponse(results, w)
}

// vehicleParseURLHandler tries to parse the URL
// encoded in one of the QR code stickers attached to the vehicle.
func vehicleParseURLHandler(w http.ResponseWriter, r *http.Request) {
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
	b, serial := adapters.ParseURL(*input.URL)
	vehicleNo := common.NormalizeVehicleNo(chi.URLParam(r, "vehicleNo"))
	if b == nil {
		serialModel := models.SerialModel{
			BureauCode: "?",
			VehicleNo:  vehicleNo,
			SerialNo:   *input.URL,
		}
		serialModel.Add()
		log.Debug().Msgf("[%s] %v", serialModel.BureauCode, serialModel)
		return
	}

	// case 3: vehicle currently offline or vehicle number does not match
	serialModel := models.SerialModel{
		BureauCode: b.Code(),
		SerialNo:   serial,
	}
	info, err := b.Info(serial)
	if err == nil {
		serialModel.VehicleNo, err = b.VehicleNo(info)
	}
	if !common.ApproxEqualVehicleNo(serialModel.VehicleNo, vehicleNo) {
		serialModel.VehicleNo = "-" + vehicleNo + "@" + serialModel.VehicleNo
		serialModel.Add()
		log.Debug().Msgf("[%s] %v", serialModel.BureauCode, serialModel)
		return
	}

	// case 4: vehicle number matches user input
	serialModel.Add()
	serialModel.AddTrainOperationLogs(info)

	results = models.ListTrainsForSingleVehicle(serialModel.VehicleNo)
	return
}
