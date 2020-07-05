package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/arnie97/emu-log/common"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter() *chi.Mux {
	mux := chi.NewRouter()
	mux.Use(
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
		middleware.Timeout(common.RequestTimeout),
	)
	mux.Get(`/map/{stationName}`, railMapHandler)
	mux.Get(`/train/{trainNo:[GDC]\d{1,4}}`, singleTrainNoHandler)
	mux.Get(`/train/{trainNo:.*,.*}`, multiTrainNoHandler)
	mux.Get(`/emu/{vehicleNo:[A-Z-0-9+]*?\d{4}}`, singleVehicleNoHandler)
	mux.Get(`/emu/{vehicleNo:[A-Z-0-9+]*?\+\d\d}`, singleVehicleNoHandler)
	mux.Get(`/emu/{vehicleNo:[A-Z-0-9+]+}`, multiVehicleNoHandler)
	return mux
}

func serializeLogEntries(rows *sql.Rows, w http.ResponseWriter) {
	results := make([]common.LogEntry, 0)
	for rows.Next() {
		var e common.LogEntry
		common.Must(rows.Scan(&e.Date, &e.VehicleNo, &e.TrainNo))
		results = append(results, e)
	}
	w.Header().Set("Content-Type", "application/json")
	common.Must(json.NewEncoder(w).Encode(results))
}
