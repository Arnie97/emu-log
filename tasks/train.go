package tasks

import (
	"strings"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/arnie97/emu-log/models"
	"github.com/rs/zerolog/log"
)

// scanTrainNo iterates over all the known vehicles operated by the specified
// railway company to see if any of these vehicles is currently associated to
// a train number (or a bunch of train numbers).
func scanTrainNo(b adapters.Bureau) {
	log.Info().Msgf("[%s] retrieving latest activities for known vehicles", b.Code())
	defer wg.Done()
	for _, serialModel := range models.ListLatestSerialForMultiVehicles(b) {
		if !strings.HasPrefix(serialModel.VehicleNo, "CR") {
			log.Debug().Msgf("[%s] %v -> ignored", b.Code(), serialModel)
			continue
		}
		info, err := b.Info(serialModel.SerialNo)

		var logModel models.LogModel
		if err == nil {
			logModel.TrainNo, logModel.Date, err = b.TrainNo(info)
		}
		if err != nil || logModel.TrainNo == "" {
			log.Debug().Msgf("[%s] %v -> %v", b.Code(), serialModel, err)
			continue
		}

		logModel.VehicleNo, err = b.VehicleNo(info)
		if common.ApproxEqualVehicleNo(serialModel.VehicleNo, logModel.VehicleNo) {
			log.Debug().Msgf("[%s] %v -> %v", b.Code(), serialModel, logModel)
			logModel.VehicleNo = serialModel.VehicleNo
			logModel.Add()
		} else {
			log.Warn().Msgf("[%s] %v -> %v ignored", b.Code(), serialModel, logModel)
			continue
		}
	}
	log.Info().Msgf("[%s] updates done for known vehicles", b.Code())
}
