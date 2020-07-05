package tasks

import (
	"database/sql"
	"time"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

// scanTrainNo iterates over all the known vehicles operated by the specified
// railway company to see if any of these vehicles is currently associated to
// a train number (or a bunch of train numbers).
func scanTrainNo(b adapters.Bureau, tx *sql.Tx) {
	log.Info().Msgf("[%s] job started: %s", b.Code(), b.Name())
	defer wg.Done()

	rows, err := tx.Query(`
		SELECT emu_no, emu_qrcode, MAX(rowid)
		FROM emu_qrcode
		WHERE emu_bureau = ?
		GROUP BY emu_no
		ORDER BY emu_no ASC;
	`, b.Code())
	common.Must(err)
	defer rows.Close()

	for rows.Next() {
		var e common.LogEntry
		var qrCode, id string
		common.Must(rows.Scan(&e.VehicleNo, &qrCode, &id))
		time.Sleep(requestDelay)
		e.TrainNo, e.Date, err = b.TrainNo(qrCode)
		if err != nil {
			log.Error().Msg(err.Error())
		}
		log.Debug().Msgf("[%s] %s -> %s", b.Code(), e.VehicleNo, e.TrainNo)
		if e.TrainNo != "" {
			res, err := tx.Exec(
				`INSERT OR IGNORE INTO emu_log VALUES (?, ?, ?)`,
				e.Date, e.VehicleNo, e.TrainNo,
			)
			common.Must(err)
			logID, err := res.LastInsertId()
			common.Must(err)

			for _, singleTrainNo := range common.NormalizeTrainNo(e.TrainNo) {
				_, err = tx.Exec(
					`REPLACE INTO emu_latest VALUES (?, ?, ?, ?)`,
					e.Date, e.VehicleNo, singleTrainNo, logID,
				)
				common.Must(err)
			}
		}
	}
	log.Info().Msgf("[%s] job done: %s", b.Code(), b.Name())
}
