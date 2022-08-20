// Package tasks defines web scraping schedules.
package tasks

import (
	"net/http"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/arnie97/emu-log/handlers"
	"github.com/rs/zerolog/log"
)

const helpMsg = `%s

usage:

%[2]s i[nfo]      SITE_ADAPTER [QR_CODE ...]
%[2]s u[nitNo]   [SITE_ADAPTER ...]
%[2]s t[rainNo]  [SITE_ADAPTER[.OPERATORS] ...]
%[2]s s[chedule] [SITE_ADAPTER[.OPERATORS] ...]
%[2]s d[aemon]
`

func CmdParser(args ...string) {
	if len(args) < 2 {
		log.Fatal().Msgf(helpMsg, "missing argument: TASK_TYPE", args[0])
	}

	switch args[1] {
	case "d", "daemon":
		serveHTTP()
	case "s", "schedule":
		go serveHTTP()
		scheduleTask(func() {
			iterateAdapters(scanTask, args[2:]...)
		})
	case "t", "trainNo":
		iterateAdapters(scanTrainNo, args[2:]...)
	case "u", "unitNo":
		iterateAdapters(scanUnitNo, args[2:]...)
	case "i", "info", "a", "add":
		if len(args) < 3 {
			log.Fatal().Msg("missing argument: SITE_ADAPTER [QR_CODE ...]")
		}

		a := adapters.MustGetAdapterByCode(args[2])
		for _, qrCode := range args[3:] {
			if args[1][0] == 'i' {
				info, err := a.Info(qrCode)
				common.PrettyPrint(info)
				common.Must(err)
			} else {
				addUnitBySerial(a, qrCode)
			}
		}

	default:
		log.Fatal().Msgf(helpMsg, "invalid TASK_TYPE: "+args[1], args[0])
	}
}

func serveHTTP() {
	const host = "localhost:8080"
	log.Info().Msgf("server listening on %s", host)
	common.Must(http.ListenAndServe(host, handlers.NewRouter()))
}
