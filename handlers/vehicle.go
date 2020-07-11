package handlers

import (
	"net/http"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/go-chi/chi"
)

type urlWrapper struct {
	URL *string `json:"url,omitempty"`
}

// singleVehicleNoHandler takes an exact vehicle number,
// and returns the 30 most recent log items for the vehicle.
func singleVehicleNoHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := common.DB().Query(`
		SELECT *
		FROM (
			SELECT *
			FROM emu_log
			WHERE emu_no IN (
				SELECT emu_no
				FROM emu_qrcode
				WHERE emu_no LIKE ?
			)
			ORDER BY date DESC, rowid DESC
			LIMIT 30
		)
		ORDER BY emu_no ASC;
	`, "%"+common.NormalizeVehicleNo(chi.URLParam(r, "vehicleNo")))
	common.Must(err)
	defer rows.Close()
	serializeLogEntries(rows, w)
}

// multiVehicleNoHandler takes an incomplete part of the vehicle number,
// and returns the most recent occurance for the first 30 vehicles
// in lexicographical order that matches the given fuzzy pattern.
func multiVehicleNoHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := common.DB().Query(`
		SELECT *
		FROM emu_log
		WHERE rowid IN (
			SELECT MAX(rowid)
			FROM emu_log
			WHERE emu_no IN (
				SELECT emu_no
				FROM emu_qrcode
				WHERE emu_no LIKE ?
			)
			GROUP BY emu_no
			LIMIT 30
		)
		ORDER BY emu_no ASC;
	`, "%"+common.NormalizeVehicleNo(chi.URLParam(r, "vehicleNo"))+"%")
	common.Must(err)
	defer rows.Close()
	serializeLogEntries(rows, w)
}

// vehicleBuildURLHandler takes an exact vehicle number, and rebuild
// the URL encoded in one of the QR code stickers attached to the vehicle.
func vehicleBuildURLHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := common.DB().Query(`
		SELECT emu_bureau, emu_qrcode
		FROM emu_qrcode
		WHERE emu_no = ?
		ORDER BY rowid DESC
		LIMIT 1;
	`, common.NormalizeVehicleNo(chi.URLParam(r, "vehicleNo")))
	common.Must(err)
	defer rows.Close()

	var results urlWrapper
	if rows.Next() {
		var bureauCode, serial string
		rows.Scan(&bureauCode, &serial)
		if b := adapters.Bureaus[bureauCode]; b != nil {
			url := adapters.BuildURL(b, serial)
			results.URL = &url
		}
	}
	jsonResponse(results, w)
}
