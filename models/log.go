package models

import (
	"strings"
	"time"

	"github.com/arnie97/emu-log/common"
)

type LogModel struct {
	Date      string `json:"date"`
	VehicleNo string `json:"emu_no"`
	TrainNo   string `json:"train_no"`
}

func (LogModel) Schema() string {
	return `
		CREATE TABLE IF NOT EXISTS emu_log (
			date        VARCHAR NOT NULL,
			emu_no      VARCHAR NOT NULL,
			train_no    VARCHAR NOT NULL,
			UNIQUE(date, emu_no, train_no)
		);
		CREATE INDEX IF NOT EXISTS idx_emu_no ON emu_log(emu_no);
		CREATE INDEX IF NOT EXISTS idx_train_no ON emu_log(train_no);

		CREATE TABLE IF NOT EXISTS emu_latest (
			date        VARCHAR NOT NULL,
			emu_no      VARCHAR NOT NULL,
			train_no    VARCHAR NOT NULL,
			log_id      INTEGER NOT NULL,
			UNIQUE(train_no)
		);
	`
}

func init() {
	Register(LogModel{})
}

// Add saves the log entry to the database,
// and update related records in the materialized view.
func (logModel LogModel) Add() {
	// use current date as the default value if date is not provided
	if logModel.Date == "" {
		logModel.Date = time.Now().Format(common.ISODate)
	}

	res, err := DB().Exec(
		`INSERT OR IGNORE INTO emu_log VALUES (?, ?, ?)`,
		logModel.Date, logModel.VehicleNo, logModel.TrainNo,
	)
	common.Must(err)

	affected, err := res.RowsAffected()
	common.Must(err)
	if affected == 0 {
		return
	}

	logID, err := res.LastInsertId()
	common.Must(err)
	logModel.Materialize(logID)
}

// Materialize updates the materialized view,
// which stores the last used vehicle for each half of train numbers.
func (logModel LogModel) Materialize(logID int64) {
	for _, singleTrainNo := range common.NormalizeTrainNo(logModel.TrainNo) {
		_, err := DB().Exec(
			`REPLACE INTO emu_latest VALUES (?, ?, ?, ?)`,
			logModel.Date, logModel.VehicleNo, singleTrainNo, logID,
		)
		common.Must(err)
	}
}

// Query executes a SQL statement and returns all the result rows.
func (logModel LogModel) Query(sql string, args ...interface{}) (logs []LogModel) {
	logs = make([]LogModel, 0)
	rows, err := DB().Query(sql, args...)
	common.Must(err)
	for rows.Next() {
		common.Must(rows.Scan(
			&logModel.Date,
			&logModel.VehicleNo,
			&logModel.TrainNo,
		))
		logs = append(logs, logModel)
	}
	return logs
}

func ListVehiclesForSingleTrain(trainNo string) []LogModel {
	return LogModel{}.Query(`
		SELECT z.date, z.emu_no, z.train_no
		FROM emu_latest AS x
		INNER JOIN emu_log AS y
		INNER JOIN emu_log AS z
		ON x.log_id = y.rowid AND y.train_no = z.train_no
		WHERE x.train_no = ?
		ORDER BY z.date DESC
		LIMIT 30;
	`, trainNo)
}

func ListLatestVehicleForMultiTrains(trainNoList []string) []LogModel {
	trainNoArgs := make([]interface{}, len(trainNoList))
	trainNoArgsPlaceHolder := strings.Repeat(", ?", len(trainNoList))[2:]
	for i := range trainNoList {
		trainNoArgs[i] = trainNoList[i]
	}
	return LogModel{}.Query(`
		SELECT date, emu_no, train_no
		FROM emu_latest
		WHERE train_no IN (`+trainNoArgsPlaceHolder+`)
	`, trainNoArgs...)
}

func ListTrainsForSingleVehicle(vehicleNo string) []LogModel {
	return LogModel{}.Query(`
		SELECT *
		FROM (
			SELECT date, emu_no, train_no
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
	`, vehicleNo)
}

func ListLatestTrainForMatchedVehicles(vehicleNo string) []LogModel {
	return ListLatestTrainByCondition(
		`LIKE ?`, "%"+common.NormalizeVehicleNo(vehicleNo)+"%")
}

func ListLatestTrainForMultiVehicles(vehicleNoList []string) []LogModel {
	vehicleNoArgs := make([]interface{}, len(vehicleNoList))
	vehicleNoPlaceHolder := strings.Repeat(", ?", len(vehicleNoList))[2:]
	for i := range vehicleNoList {
		vehicleNoArgs[i] = common.NormalizeVehicleNo(vehicleNoList[i])
	}
	return ListLatestTrainByCondition(
		`IN (`+vehicleNoPlaceHolder+`)`, vehicleNoArgs...)
}

func ListLatestTrainByCondition(cond string, args ...interface{}) []LogModel {
	return LogModel{}.Query(`
		SELECT date, emu_no, train_no
		FROM emu_log
		WHERE rowid IN (
			SELECT MAX(rowid)
			FROM emu_log
			WHERE emu_no IN (
				SELECT DISTINCT emu_no
				FROM emu_qrcode
				WHERE emu_no `+cond+`
			)
			GROUP BY emu_no
			LIMIT 30
		)
		ORDER BY emu_no
	`, args...)
}
