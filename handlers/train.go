package handlers

import (
	"net/http"
	"strings"

	"github.com/arnie97/emu-log/models"
	"github.com/go-chi/chi"
)

// singleTrainNoHandler returns the used vehicle and the corresponding date
// for the 30 most recent log items that matches the given train number.
func singleTrainNoHandler(w http.ResponseWriter, r *http.Request) {
	trainNo := chi.URLParam(r, "trainNo")
	results := models.ListVehiclesForSingleTrain(trainNo)
	jsonResponse(results, w)
}

// multiTrainNoHandler returns the last used vehicle for multiple trains.
func multiTrainNoHandler(w http.ResponseWriter, r *http.Request) {
	trainNoList := strings.Split(chi.URLParam(r, "trainNo"), ",")
	results := models.ListLatestVehicleForMultiTrains(trainNoList)
	jsonResponse(results, w)
}
