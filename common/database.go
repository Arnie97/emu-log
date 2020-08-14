package common

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/mattn/go-sqlite3"
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
	CREATE INDEX IF NOT EXISTS idx_train_no ON emu_log(train_no);
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
		var err error
		dbConn, err = sql.Open("sqlite3", AppPath()+"/emu-log.db")
		Must(err)
		// TODO: defer dbConn.Close()

		_, err = dbConn.Exec(dbSchema)
		Must(err)
	})

	return dbConn
}

// CountRecords takes a tableName and returns the number of rows in the table.
func CountRecords(tableName string, fields ...string) (count int) {
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
