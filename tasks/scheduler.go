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
	sched := common.Conf().Schedule
	startTime := time.Duration(sched.StartTime)
	endTime := time.Duration(sched.EndTime)

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
		if !nextRun.IsZero() {
			log.Info().Msgf("next scheduled run: %v", nextRun)
			time.Sleep(time.Until(nextRun))
		}
		task()
	}
}

// scanTask is a combination of scanUnitNo() and scanTrainNo().
func scanTask(a adapters.Adapter) {
	endTime := time.Duration(common.Conf().Schedule.EndTime)
	scanForNewUnits := true
	if a.AlwaysOn() {
		now := time.Now()
		today := time.Date(
			now.Year(), now.Month(), now.Day(),
			0, 0, 0, 0, time.Local,
		)

		// for "always on" adapters, it would be more than sufficient to scan
		// the whole key space once a day to discover recently added units.
		// let's run it during the possessive intervals in the train diagrams.
		if !now.After(today.Add(endTime)) {
			scanForNewUnits = false
		}
	}

	if scanForNewUnits {
		wg.Add(1)
		defer scanUnitNo(a)
	}
	scanTrainNo(a)
}

// iterateAdapters parallelizes scanning requests for different railway
// companies with goroutines.
func iterateAdapters(task func(adapters.Adapter), adapterList ...string) {
	once.Do(func() {
		checkLocalTimezone()
		checkInternetConnection()
		checkDatabase()
	})

	// support both joined adapter codes and space separated adapter codes
	adapterCodes := strings.Join(adapterList, "")
	if len(adapterCodes) == 0 {
		for _, a := range adapters.Adapters {
			wg.Add(1)
			go task(a)
		}
	} else {
		for _, code := range adapterCodes {
			a := adapters.MustGetAdapterByCode(string(code))
			wg.Add(1)
			go task(a)
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
	_, err := adapters.Adapters["K"].Info("K1001036127001")
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
		"found %d units and %d QR codes in the database",
		models.CountRecords("emu_qr_code", "DISTINCT emu_no"),
		models.CountRecords("emu_qr_code"),
	)
}
