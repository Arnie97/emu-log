package tasks

import (
	"database/sql"
	"time"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

// scanVehicleNo trys each unknown QR code in the brute force key space to see
// if any of these serial numbers was recently (or is currently) put in to use.
func scanVehicleNo(b adapters.Bureau, tx *sql.Tx) {
	log.Info().Msgf("[%s] started scanning for new vehicles", b.Code())
	defer wg.Done()

	rows, err := tx.Query(`
		SELECT emu_qrcode
		FROM emu_qrcode
		WHERE emu_bureau = ?
		ORDER BY emu_qrcode ASC;
	`, b.Code())
	common.Must(err)
	defer rows.Close()

	serials := make(chan string)
	go func() {
		b.BruteForce(serials)
		close(serials)
	}()

	var serialFromDB string
	for serial := range serials {
		// skip existing codes in the database
		for serial > serialFromDB && rows.Next() {
			common.Must(rows.Scan(&serialFromDB))
		}
		if serial == serialFromDB {
			continue
		}
		time.Sleep(requestDelay)
		addVehicleBySerial(b, tx, serial)
	}
	log.Info().Msgf("[%s] finished scanning", b.Code())
}

// addVehicleBySerial takes a serial number from some railway company and
// save it to the database if the serial number maps to a vehicle number.
func addVehicleBySerial(b adapters.Bureau, tx *sql.Tx, serial string) {
	// handle errors
	var e common.LogEntry
	info, err := b.Info(serial)
	if err == nil {
		e.VehicleNo, err = b.VehicleNo(info)
	}
	if err != nil || e.VehicleNo == "" {
		log.Debug().Msgf("[%s] %s -> %v", b.Code(), serial, err)
		return
	}

	// add a vehicle serial number record
	_, err = tx.Exec(
		`INSERT OR IGNORE INTO emu_qrcode VALUES (?, ?, ?)`,
		e.VehicleNo, b.Code(), serial,
	)
	common.Must(err)

	// also add a activity log record if the train number is available
	e.TrainNo, e.Date, err = b.TrainNo(info)
	if err == nil && e.TrainNo != "" {
		addTrainOperationLog(&e, tx)
	}
	log.Debug().Msgf("[%s] %s -> %v", b.Code(), serial, e)
}
