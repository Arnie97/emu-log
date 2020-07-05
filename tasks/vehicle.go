package tasks

import (
	"database/sql"
	"time"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

func scanVehicleNo(b adapters.Bureau, tx *sql.Tx) {
	log.Info().Msgf("[%s] job started: %s", b.Code(), b.Name())
	defer wg.Done()

	rows, err := tx.Query(`
		SELECT emu_qrcode
		FROM emu_qrcode
		WHERE emu_bureau = ?
		ORDER BY emu_qrcode ASC;
	`, b.Code())
	common.Must(err)
	defer rows.Close()

	qrCodes := make(chan string)
	go func() {
		b.BruteForce(qrCodes)
		close(qrCodes)
	}()

	qrCodeFromDB := ""
	for qrCode := range qrCodes {
		// skip existing codes in the database
		for qrCode > qrCodeFromDB && rows.Next() {
			common.Must(rows.Scan(&qrCodeFromDB))
			log.Debug().Msgf("[%s] loaded: %s", b.Code(), qrCodeFromDB)
		}
		if qrCode == qrCodeFromDB {
			continue
		}

		time.Sleep(requestDelay)
		vehicleNo, err := b.VehicleNo(qrCode)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		log.Debug().Msgf("[%s] checked: %s -> %s", b.Code(), qrCode, vehicleNo)
		if vehicleNo != "" {
			_, err := tx.Exec(
				`INSERT OR IGNORE INTO emu_qrcode VALUES (?, ?, ?)`,
				vehicleNo, b.Code(), qrCode,
			)
			common.Must(err)
		}
	}
	log.Info().Msgf("[%s] job done: %s", b.Code(), b.Name())
}
