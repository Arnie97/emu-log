package handlers

import (
	"net/http"
	"strings"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/arnie97/emu-log/models"
	"github.com/go-chi/chi"
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
