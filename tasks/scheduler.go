package tasks

import (
	"strings"
	"sync"
	"time"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/arnie97/emu-log/models"
	"github.com/rs/zerolog/log"
)

const (
	startTime = 6 * time.Hour
	endTime   = 22 * time.Hour
	day       = 24 * time.Hour
	tzBeijing = 8 * time.Hour
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
		} else if now.After(today.Add(endTime)) {
			nextRun = today.Add(startTime + day)
		}
		task()
		if !nextRun.IsZero() {
			log.Info().Msgf("next scheduled run: %v", nextRun)
			time.Sleep(time.Until(nextRun))
		}
	}
}

// scanTask is a combination of scanVehicleNo() and scanTrainNo().
func scanTask(b adapters.Bureau) {
	scanForNewVehicles := true
	if b.AlwaysOn() {
		now := time.Now()
		today := time.Date(
			now.Year(), now.Month(), now.Day(),
			0, 0, 0, 0, time.Local,
		)

		// for "always on" adapters, it would be more than sufficient to scan
		// the whole key space once a day to discover recently added vehicles.
		// let's run it during the possessive intervals in the train diagrams.
		if !now.After(today.Add(endTime)) {
			scanForNewVehicles = false
		}
	}

	if scanForNewVehicles {
		wg.Add(1)
		defer scanVehicleNo(b)
	}
	scanTrainNo(b)
}

// iterateBureaus parallelizes scanning requests for different railway
// companies with goroutines.
func iterateBureaus(task func(adapters.Bureau), bureaus ...string) {
	once.Do(func() {
		checkLocalTimezone()
		checkInternetConnection()
		checkDatabase()
	})

	// support both joined bureau codes and space separated bureau codes
	bureauCodes := strings.Join(bureaus, "")
	if len(bureauCodes) == 0 {
		for _, b := range adapters.Bureaus {
			wg.Add(1)
			go task(b)
		}
	} else {
		for _, code := range bureauCodes {
			b := adapters.MustGetBureauByCode(string(code))
			wg.Add(1)
			go task(b)
		}
	}

	wg.Wait()
}

// checkLocalTimezone prints a warning if the server timezone settings is
// different from China Railways (UTC+08).
func checkLocalTimezone() {
	tzName, tzOffset := time.Now().Zone()
	if time.Duration(tzOffset)*time.Second != tzBeijing {
		log.Warn().Msgf(
			"expected Beijing Timezone (UTC+08), but found %s (UTC%s)",
			tzName, time.Now().Format("-07"),
		)
	}
}

// checkInternetConnection prints the RTT for a HTTP connection.
func checkInternetConnection() {
	start := time.Now()
	_, err := adapters.Bureaus["H"].Info("PQ1234567")
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
		models.CountRecords("emu_log"),
	)
	log.Info().Msgf(
		"found %d vehicles and %d qr codes in the database",
		models.CountRecords("emu_qrcode", "DISTINCT emu_no"),
		models.CountRecords("emu_qrcode"),
	)
}
