package models

import (
	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

type SerialModel struct {
	VehicleNo  string `json:"emu_no"`
	BureauCode string `json:"emu_bureau"`
	SerialNo   string `json:"emu_qrcode"`
}

func (SerialModel) Schema() string {
	return `
		CREATE TABLE IF NOT EXISTS emu_qrcode (
			emu_no      VARCHAR NOT NULL,
			emu_bureau  CHAR(1) NOT NULL,
			emu_qrcode  VARCHAR NOT NULL,
			UNIQUE(emu_bureau, emu_qrcode)
		);
	`
}

func init() {
	Register(SerialModel{})
}

// Add inserts a recently discovered serial number into the database.
func (s SerialModel) Add() {
	_, err := DB().Exec(
		`INSERT OR IGNORE INTO emu_qrcode VALUES (?, ?, ?)`,
		s.VehicleNo, s.BureauCode, s.SerialNo,
	)
	common.Must(err)
}

// AddTrainOperationLogs creates related operation log records if possible.
func (s SerialModel) AddTrainOperationLogs(info adapters.JSONObject) {
	b := adapters.MustGetBureauByCode(s.BureauCode)
	trains, err := b.TrainNo(info)
	if err != nil || len(trains) == 0 {
		log.Debug().Msgf("[%s] %v -> %v", b.Code(), s, err)
		return
	}

	var logModel LogModel
	logModel.VehicleNo, _ = b.VehicleNo(info)
	for _, train := range trains {
		logModel.TrainNo = train.TrainNo
		logModel.Date = train.Date
		if !common.ApproxEqualVehicleNo(s.VehicleNo, logModel.VehicleNo) {
			log.Warn().Msgf("[%s] %v -> %v ignored", b.Code(), s, logModel)
			return
		}
		log.Debug().Msgf("[%s] %v -> %v", b.Code(), s, logModel)
		logModel.VehicleNo = s.VehicleNo
		logModel.Add()
	}
}

// Query executes a SQL statement and returns all the result rows.
func (s SerialModel) Query(sql string, args ...interface{}) (serials []SerialModel) {
	rows, err := DB().Query(sql, args...)
	common.Must(err)
	defer rows.Close()

	for rows.Next() {
		common.Must(rows.Scan(&s.VehicleNo, &s.BureauCode, &s.SerialNo))
		serials = append(serials, s)
	}
	return serials
}

// ListSerials returns all known serials and corresponding vehicles
// of the given railway company from the database.
func ListSerials(b adapters.Bureau) []SerialModel {
	return SerialModel{}.Query(`
		SELECT emu_no, emu_bureau, emu_qrcode
		FROM emu_qrcode
		WHERE emu_bureau = ?
		ORDER BY emu_qrcode ASC;
	`, b.Code())
}

// ListSerialsForSingleVehicle returns all the known serials for a vehicle.
func ListSerialsForSingleVehicle(vehicleNo string) []SerialModel {
	return SerialModel{}.Query(`
		SELECT emu_no, emu_bureau, emu_qrcode
		FROM emu_qrcode
		WHERE emu_no LIKE ?
		ORDER BY rowid DESC;
	`, vehicleNo)
}

// ListLatestSerialForMultiVehicles returns the most recently discovered
// serial number for each vehicle in the given railway company, but excluding
// those with known train schedules.
func ListLatestSerialForMultiVehicles(b adapters.Bureau) []SerialModel {
	return SerialModel{}.Query(`
		SELECT emu_qrcode.emu_no, emu_bureau, emu_qrcode
		FROM (
			SELECT emu_no, emu_bureau, emu_qrcode
			FROM emu_qrcode
			WHERE emu_bureau = ?
			GROUP BY emu_no
			HAVING MAX(rowid)
			ORDER BY emu_no ASC
		) AS emu_qrcode
		LEFT JOIN (
			SELECT emu_no, date
			FROM (
				SELECT rowid, emu_no, date
				FROM emu_log
				ORDER BY rowid DESC
				LIMIT 5000
			)
			GROUP BY emu_no
			HAVING MAX(rowid)
		) AS emu_log
		ON emu_qrcode.emu_no = emu_log.emu_no
		WHERE date < DATETIME('now', 'localtime');
	`, b.Code())
}
