package tasks

import (
	"net/http"
	"os"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/arnie97/emu-log/handlers"
	"github.com/rs/zerolog/log"
)

const helpMsg = `missing argument: TASK_TYPE

usage:

%[1]s info       BUREAU_CODE [QR_CODE ...]
%[1]s trainNo   [BUREAU_CODE ...]
%[1]s vehicleNo [BUREAU_CODE ...]
%[1]s schedule  [BUREAU_CODE ...]
%[1]s serve
`

func CmdParser() {
	if len(os.Args) < 2 {
		log.Fatal().Msgf(helpMsg, os.Args[0])
	}

	switch os.Args[1] {
	case "serve":
		serveHTTP()
	case "schedule":
		go serveHTTP()
		scheduleTask(func() {
			iterateBureaus(task, os.Args[2:]...)
		})
	case "trainNo":
		iterateBureaus(scanTrainNo, os.Args[2:]...)
	case "vehicleNo":
		iterateBureaus(scanVehicleNo, os.Args[2:]...)
	case "info":
		if len(os.Args) < 3 {
			log.Fatal().Msg("missing argument: BUREAU_CODE [QR_CODE ...]")
		}

		b := adapters.MustGetBureauByCode(os.Args[2])
		for _, qrCode := range os.Args[3:] {
			info, _ := b.Info(qrCode)
			common.PrettyPrint(info)
		}
	default:
		log.Fatal().Msgf("invalid TASK_TYPE: %s", os.Args[1])
	}
}

func serveHTTP() {
	const host = "localhost:8080"
	log.Info().Msgf("server listening on %s", host)
	common.Must(http.ListenAndServe(host, handlers.NewRouter()))
}
