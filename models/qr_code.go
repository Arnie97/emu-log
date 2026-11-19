package models

import "github.com/arnie97/emu-log/common"

type QRCodeModel struct {
	Adapter  string  `json:"adapter"`
	Operator string  `json:"operator"`
	UnitNo   string  `json:"emu_no"`
	QRCode   string  `json:"qr_code"`
	Date     *string `json:"date,omitempty"`
}

func (q QRCodeModel) String() string {
	return formatQRCodeRecord(q.Adapter, q.Operator, q.UnitNo, q.QRCode, q.Date)
}

func (QRCodeModel) Schema() string {
	return `
		CREATE TABLE IF NOT EXISTS emu_qr_code (
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
	Register(QRCodeModel{})
}

// Add inserts a QR code into the database.
func (q QRCodeModel) Add() {
	_, err := DB().Exec(
		`INSERT OR IGNORE INTO
		emu_qr_code(adapter, operator, emu_no, qr_code, date) VALUES (?, ?, ?, ?, ?)`,
		q.Adapter, q.Operator, q.UnitNo, q.QRCode, common.DateTimeString(),
	)
	common.Must(err)
}
