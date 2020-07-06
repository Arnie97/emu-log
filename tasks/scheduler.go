package tasks

import (
	"database/sql"
	"strings"
	"sync"
	"time"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/rs/zerolog/log"
)

const (
	day            = 24 * time.Hour
	repeatInterval = time.Hour
	requestDelay   = 3 * time.Second
	startTime      = 5 * time.Hour
	endTime        = 24 * time.Hour
	beijingTime    = 8 * time.Hour
)

var (
	wg   sync.WaitGroup
	once sync.Once
)

// scheduleTask run web scraping routines on an regular basis.
func scheduleTask(task func()) {
	var nextRun time.Time
	for {
		now := time.Now()
		today := time.Date(
			now.Year(), now.Month(), now.Day(),
			0, 0, 0, 0, time.Local,
		)
		if now.Before(today.Add(startTime)) {
			nextRun = today.Add(startTime)
		} else if now.After(today.Add(endTime - repeatInterval)) {
			nextRun = today.Add(startTime + day)
		} else {
			nextRun = now.Truncate(repeatInterval).Add(repeatInterval)
		}
		log.Info().Msgf("next scheduled run: %v", nextRun)
		time.Sleep(time.Until(nextRun))
		task()
	}
}

// scanTask is a combination of scanVehicleNo() and scanTrainNo().
func scanTask(b adapters.Bureau, tx *sql.Tx) {
	scanForNewVehicles := false

	// these bureau adapters return nothing when online ordering is disabled,
	// so we cannot distinguish nonexistent barcodes from offline vehicles,
	// and should always scan the whole key space.
	const bureausAlwaysScanAll = "PQ"
	if strings.Contains(bureausAlwaysScanAll, b.Code()) {
		scanForNewVehicles = true
	} else {
		now := time.Now()
		today := time.Date(
			now.Year(), now.Month(), now.Day(),
			0, 0, 0, 0, time.Local,
		)

		// for other bureau adapters, it would be more than sufficient to scan
		// the whole key space once a day to discover recently added vehicles.
		// let's run it during the possessive intervals in the train diagrams.
		if now.After(today.Add(endTime - repeatInterval)) {
			scanForNewVehicles = true
		}
	}

	if scanForNewVehicles {
		wg.Add(1)
		defer scanVehicleNo(b, tx)
	}
	scanTrainNo(b, tx)
}

// iterateBureaus parallelizes scanning requests for different railway
// companies with goroutines.
func iterateBureaus(task func(adapters.Bureau, *sql.Tx), bureaus ...string) {
	once.Do(func() {
		checkLocalTimezone()
		checkInternetConnection()
	})

	tx, err := common.DB().Begin()
	common.Must(err)
	defer tx.Rollback()

	// support both joined bureau codes and space separated bureau codes
	bureauCodes := strings.Join(bureaus, "")
	if len(bureauCodes) == 0 {
		for _, b := range adapters.Bureaus {
			wg.Add(1)
			go task(b, tx)
		}
	} else {
		for _, code := range bureauCodes {
			b := adapters.MustGetBureauByCode(string(code))
			wg.Add(1)
			go task(b, tx)
		}
	}

	wg.Wait()
	tx.Commit()
}

// checkLocalTimezone prints a warning if the server timezone settings is
// different from China Railways (UTC+08).
func checkLocalTimezone() {
	tzName, tzOffset := time.Now().Zone()
	if time.Duration(tzOffset)*time.Second != beijingTime {
		log.Warn().Msgf(
			"expected Beijing Timezone (UTC+08), but found %s (UTC%s)",
			tzName, time.Now().Format("-07"),
		)
	}
}

// checkInternetConnection prints the RTT for a HTTP connection.
func checkInternetConnection() {
	start := time.Now()
	_, err := adapters.Bureaus["H"].Info("PQ0123456")
	common.Must(err)
	log.Info().Msgf(
		"internet connection ok, round-trip delay %v",
		time.Since(start),
	)
}

// checkDatabase prints row counts for all tables to ensure a working DB
func checkDatabase() {
	log.Info().Msgf(
		"found %d log records in the database",
		common.CountRecords("emu_log"),
	)
	log.Info().Msgf(
		"found %d vehicles and %d qr codes in the database",
		common.CountRecords("emu_qrcode", "DISTINCT emu_no"),
		common.CountRecords("emu_qrcode"),
	)
}
