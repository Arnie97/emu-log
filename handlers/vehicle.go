package handlers

import (
	"net/http"

	"github.com/arnie97/emu-log/common"
	"github.com/go-chi/chi"
)

// singleVehicleNoHandler returns the 30 most recent log items for the given
// vehicle number.
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

// multiVehicleNoHandler returns the most recent occurance for the first 30
// vehicles in lexicographical order that matches the given fuzzy pattern.
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
