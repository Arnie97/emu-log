package tasks

import (
	"strings"

	"github.com/arnie97/emu-log/adapters"
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
		if err != nil {
			log.Debug().Msgf("[%s] %v -> %v", b.Code(), serialModel, err)
		}
		serialModel.AddTrainOperationLogs(info)
	}
	log.Info().Msgf("[%s] updates done for known vehicles", b.Code())
}
