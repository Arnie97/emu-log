package tasks

import (
	"strings"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/models"
	"github.com/rs/zerolog/log"
)

// scanTrainNo iterates over all the known units operated by the specified
// adapter to see if any of these units is currently associated to a train
// number (or a bunch of train numbers).
func scanTrainNo(a adapters.Adapter) {
	log.Info().Msgf("[%s] retrieving latest activities for known units", a.Code())
	defer wg.Done()
	for _, serialModel := range models.ListLatestSerialForMultiUnits(a) {
		if !strings.HasPrefix(serialModel.UnitNo, "CR") {
			log.Debug().Msgf("[%s] %v -> ignored", a.Code(), serialModel)
			continue
		}
		info, err := a.Info(serialModel.SerialNo)
		if err != nil {
			log.Debug().Msgf("[%s] %v -> %v", a.Code(), serialModel, err)
		}
		serialModel.AddTrainOperationLogs(info)
	}
	log.Info().Msgf("[%s] updates done for known units", a.Code())
}
