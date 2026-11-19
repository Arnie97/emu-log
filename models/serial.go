package models

import (
	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

type SerialModel struct {
	Adapter  string  `json:"adapter"`
	Operator string  `json:"operator"`
	UnitNo   string  `json:"emu_no"`
	SerialNo string  `json:"qr_code"`
	Date     *string `json:"date,omitempty"`
}

func (s SerialModel) String() string {
	return formatQRCodeRecord(s.Adapter, s.Operator, s.UnitNo, s.SerialNo, s.Date)
}

func (SerialModel) Schema() string {
	return `
		CREATE TABLE IF NOT EXISTS emu_serial (
			adapter  CHAR(1) NOT NULL,
			operator CHAR(1) NOT NULL,
			emu_no   VARCHAR NOT NULL,
			qr_code  VARCHAR NOT NULL,
			date     VARCHAR,
			UNIQUE(adapter, qr_code)
		);
	`
}

func init() {
	Register(SerialModel{})
}

// Add inserts a recently discovered serial number into the database.
func (s SerialModel) Add() {
	_, err := DB().Exec(
		`INSERT OR IGNORE INTO
		emu_serial(adapter, operator, emu_no, qr_code, date) VALUES (?, ?, ?, ?, ?)`,
		s.Adapter, s.Operator, s.UnitNo, s.SerialNo, common.DateTimeString(),
	)
	common.Must(err)
}

// AddTrainOperationLogs creates related operation log records if possible.
func (s SerialModel) AddTrainOperationLogs(info adapters.JSONObject) {
	a := adapters.MustGetAdapterByCode(s.Adapter)
	trains, err := a.TrainNo(info)
	if err != nil || len(trains) == 0 {
		log.Debug().Msgf("[%s] %v -> %v", a.Code(), s, err)
		return
	}

	var logModel LogModel
	logModel.UnitNo, _ = a.UnitNo(s.SerialNo, info)
	for _, train := range trains {
		logModel.TrainNo = train.TrainNo
		logModel.Date = train.Date
		if !common.ApproxEqualUnitNo(s.UnitNo, logModel.UnitNo) {
			log.Warn().Msgf("[%s] %v -> %v ignored", a.Code(), s, logModel)
			return
		}
		log.Debug().Msgf("[%s] %v -> %v", a.Code(), s, logModel)
		logModel.UnitNo = s.UnitNo
		logModel.Add()
	}
}

// Query executes a SQL statement and returns all the result rows.
func (s SerialModel) Query(query string, args ...interface{}) (serials []SerialModel) {
	rows, err := DB().Query(query, args...)
	common.Must(err)
	defer rows.Close()

	for rows.Next() {
		common.Must(rows.Scan(&s.Adapter, &s.Operator, &s.UnitNo, &s.SerialNo, &s.Date))
		serials = append(serials, s)
	}
	common.Must(rows.Err())
	return serials
}

// ListSerials returns all known serials and corresponding units
// of a given adapter from the database.
func ListSerials(a adapters.Adapter) []SerialModel {
	return SerialModel{}.Query(`
		SELECT adapter, operator, emu_no, qr_code, date
		FROM emu_serial
		WHERE adapter = ?
		ORDER BY qr_code ASC;
	`, a.Code())
}

// ListSerialsForSingleUnit returns all the known serials for one unit.
func ListSerialsForSingleUnit(unitNo string) []SerialModel {
	return SerialModel{}.Query(`
		SELECT adapter, operator, emu_no, qr_code, date
		FROM emu_serial
		WHERE emu_no LIKE ?
		ORDER BY rowid DESC;
	`, unitNo)
}

// ListLatestSerialForMultiUnits returns the most recently discovered
// serial number for each unit from the given adapter, but excluding
// those with known train schedules.
func ListLatestSerialForMultiUnits(a adapters.Adapter, operators ...string) []SerialModel {
	query, argv := In(`
		SELECT adapter, operator, emu_serial.emu_no, qr_code, emu_serial.date
		FROM (
			SELECT emu_no, adapter, operator, qr_code, date
			FROM emu_serial
			WHERE operator %s AND adapter = ?
			GROUP BY emu_no
			HAVING MAX(rowid)
			ORDER BY emu_no ASC
		) AS emu_serial
		LEFT JOIN (
			SELECT emu_no, date
			FROM (
				SELECT rowid, emu_no, date
				FROM emu_log
				ORDER BY rowid DESC
				LIMIT 10000
			)
			GROUP BY emu_no
			HAVING MAX(rowid)
		) AS emu_log
		ON emu_serial.emu_no = emu_log.emu_no
		WHERE emu_log.date IS NULL OR emu_log.date < DATETIME('now', 'localtime');
	`, operators, "IS NOT NULL")
	argv = append(argv, a.Code())
	return SerialModel{}.Query(query, argv...)
}
