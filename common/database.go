package common

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

const dbSchema = `
	CREATE TABLE IF NOT EXISTS emu_latest (
		date        VARCHAR NOT NULL,
		emu_no      VARCHAR NOT NULL,
		train_no    VARCHAR NOT NULL,
		log_id      INTEGER NOT NULL,
		UNIQUE(train_no)
	);
	CREATE TABLE IF NOT EXISTS emu_log (
		date        VARCHAR NOT NULL,
		emu_no      VARCHAR NOT NULL,
		train_no    VARCHAR NOT NULL,
		UNIQUE(date, emu_no, train_no)
	);
	CREATE TABLE IF NOT EXISTS emu_qrcode (
		emu_no      VARCHAR NOT NULL,
		emu_bureau  CHAR(1) NOT NULL,
		emu_qrcode  VARCHAR NOT NULL,
		UNIQUE(emu_bureau, emu_qrcode)
	);
	CREATE INDEX IF NOT EXISTS idx_emu_no ON emu_log(emu_no);
`

type (
	LogEntry struct {
		Date      string `json:"date"`
		VehicleNo string `json:"emu_no"`
		TrainNo   string `json:"train_no"`
	}
)

var (
	dbConn *sql.DB
	dbOnce sync.Once
)

func DB() *sql.DB {
	dbOnce.Do(func() {
		dbConn, err := sql.Open("sqlite3", "./emu_log.db")
		Must(err)
		// TODO: defer dbConn.Close()

		_, err = dbConn.Exec(dbSchema)
		Must(err)
		dbStatistics()
	})

	return dbConn
}

func dbStatistics() {
	log.Info().Msgf(
		"found %d log records in the database",
		countRecords("emu_log"),
	)
	log.Info().Msgf(
		"found %d vehicles and %d qr codes in the database",
		countRecords("emu_qrcode", "DISTINCT emu_no"),
		countRecords("emu_qrcode"),
	)
}

func countRecords(tableName string, fields ...string) (count int) {
	field := "*"
	if len(fields) != 0 {
		field = fields[0]
	}
	row := DB().QueryRow(fmt.Sprintf(
		`SELECT COUNT(%s) FROM %s`, field, tableName,
	))
	Must(row.Scan(&count))
	return
}
