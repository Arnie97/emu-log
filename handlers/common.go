// Package handlers defines service handlers for incoming HTTP requests.
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/arnie97/emu-log/common"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// NewRouter defines middlewares and API endpoints for the HTTP service.
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

// serializeLogEntries converts database query results to a JSON array.
func serializeLogEntries(rows *sql.Rows, w http.ResponseWriter) {
	results := make([]common.LogEntry, 0)
	for rows.Next() {
		var e common.LogEntry
		common.Must(rows.Scan(&e.Date, &e.VehicleNo, &e.TrainNo))
		results = append(results, e)
	}
	jsonResponse(results, w)
}

// jsonResponse takes a structure and marshals it to a JSON HTTP response.
func jsonResponse(v interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", common.ContentType)
	common.Must(json.NewEncoder(w).Encode(v))
}
