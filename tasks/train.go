package tasks

import (
	"database/sql"
	"strings"
	"time"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

// scanTrainNo iterates over all the known vehicles operated by the specified
// railway company to see if any of these vehicles is currently associated to
// a train number (or a bunch of train numbers).
func scanTrainNo(b adapters.Bureau, tx *sql.Tx) {
	log.Info().Msgf("[%s] retrieving latest activities for known vehicles", b.Code())
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
		var (
			e      common.LogEntry
			serial string
			id     int64
		)
		common.Must(rows.Scan(&e.VehicleNo, &serial, &id))
		if !strings.HasPrefix(e.VehicleNo, "CR") {
			log.Debug().Msgf("[%s] %s -> ignored", b.Code(), e.VehicleNo)
			continue
		}
		time.Sleep(requestDelay)
		info, err := b.Info(serial)
		if err == nil {
			e.TrainNo, e.Date, err = b.TrainNo(info)
		}
		if err != nil || e.TrainNo == "" {
			log.Debug().Msgf("[%s] %s -> %v", b.Code(), e.VehicleNo, err)
			continue
		}

		vehicleNo, err := b.VehicleNo(info)
		if vehicleNo[len(vehicleNo)-4] == e.VehicleNo[len(e.VehicleNo)-4] ||
			strings.ContainsRune(vehicleNo, '@') {
			log.Debug().Msgf("[%s] %s -> %v", b.Code(), vehicleNo, e)
			addTrainOperationLog(&e, tx)
		} else {
			log.Warn().Msgf("[%s] %s -> %v ignored", b.Code(), vehicleNo, e)
			continue
		}
	}
	log.Info().Msgf("[%s] updates done for known vehicles", b.Code())
}

// addTrainOperationLog saves the log entry to DB,
// and update related records in the materialized view.
func addTrainOperationLog(e *common.LogEntry, tx *sql.Tx) {
	// use current date as the default value if date is not provided
	if e.Date == "" {
		e.Date = time.Now().Format(common.ISODate)
	}

	res, err := tx.Exec(
		`INSERT OR IGNORE INTO emu_log VALUES (?, ?, ?)`,
		e.Date, e.VehicleNo, e.TrainNo,
	)
	common.Must(err)

	affected, err := res.RowsAffected()
	common.Must(err)
	if affected == 0 {
		return
	}

	// update the materialized view: last used vehicle for each train number
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

func rebuildLatest(tx *sql.Tx) {
	rows, err := tx.Query(`
		SELECT rowid, emu_no, train_no, date
		FROM emu_log
		ORDER BY rowid ASC;
	`)
	common.Must(err)
	defer rows.Close()

	log.Info().Msg("rebuilding materialized view")
	for rows.Next() {
		var (
			logID int64
			e     common.LogEntry
		)
		common.Must(rows.Scan(&logID, &e.VehicleNo, &e.TrainNo, &e.Date))
		for _, singleTrainNo := range common.NormalizeTrainNo(e.TrainNo) {
			_, err = tx.Exec(
				`REPLACE INTO emu_latest VALUES (?, ?, ?, ?)`,
				e.Date, e.VehicleNo, singleTrainNo, logID,
			)
			common.Must(err)
		}
	}
	log.Info().Msg("commiting changes")
}
