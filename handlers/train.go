package handlers

import (
	"net/http"
	"strings"

	"github.com/arnie97/emu-log/models"
	"github.com/go-chi/chi"
)

// singleTrainNoHandler returns the used unit and the corresponding date
// for the 30 most recent log items that matches the given train number.
func singleTrainNoHandler(w http.ResponseWriter, r *http.Request) {
	trainNo := chi.URLParam(r, "trainNo")
	results := models.ListUnitsForSingleTrainNo(trainNo)
	jsonResponse(results, w)
}

// multiTrainNoHandler returns the last used unit number for multiple trains.
func multiTrainNoHandler(w http.ResponseWriter, r *http.Request) {
	trainNoList := strings.Split(chi.URLParam(r, "trainNo"), ",")
	results := models.ListLatestUnitForMultiTrains(trainNoList)
	jsonResponse(results, w)
}
