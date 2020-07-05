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
	requestDelay   = 2 * time.Second
	startTime      = 5 * time.Hour
	endTime        = 24 * time.Hour
	beijingTime    = 8 * time.Hour
)

var (
	wg   sync.WaitGroup
	once sync.Once
)

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

func task(b adapters.Bureau, tx *sql.Tx) {
	now := time.Now()
	today := time.Date(
		now.Year(), now.Month(), now.Day(),
		0, 0, 0, 0, time.Local,
	)
	if b.Code() != "H" || now.After(today.Add(endTime-repeatInterval)) {
		wg.Add(1)
		defer scanVehicleNo(b, tx)
	}
	scanTrainNo(b, tx)
}

func iterateBureaus(task func(adapters.Bureau, *sql.Tx), bureauArray ...string) {
	once.Do(func() {
		checkLocalTimezone()
		checkInternetConnection()
	})

	tx, err := common.DB().Begin()
	common.Must(err)
	defer tx.Rollback()

	// support both joined bureau codes and space separated bureau codes
	bureauCodes := strings.Join(bureauArray, "")
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

func checkLocalTimezone() {
	tzName, tzOffset := time.Now().Zone()
	if time.Duration(tzOffset)*time.Second != beijingTime {
		log.Warn().Msgf(
			"expected Beijing Timezone (UTC+08), but found %s (UTC%s)",
			tzName, time.Now().Format("-07"),
		)
	}
}

func checkInternetConnection() {
	start := time.Now()
	_, err := adapters.Bureaus["H"].Info("PQ0123456")
	common.Must(err)
	log.Info().Msgf(
		"internet connection ok, round-trip delay %v",
		time.Since(start),
	)
}
