package tasks

import (
	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/models"
	"github.com/rs/zerolog/log"
)

// scanUnitNo tries each unknown QR code in the brute force key space to see
// if any of these serial numbers was recently (or is currently) put in to use.
func scanUnitNo(a adapters.Adapter) {
	log.Info().Msgf("[%s] started scanning for new units", a.Code())
	defer wg.Done()

	serials := make(chan string)
	go func() {
		for _, rule := range adapters.AdapterConf(a).SearchSpace {
			rule.Emit(serials)
		}
		close(serials)
	}()

	var index int
	serialModels := models.ListSerials(a)
	for serial := range serials {
		// skip existing serial numbers in the database
		for index < len(serialModels) && serialModels[index].SerialNo < serial {
			index++
		}
		if index < len(serialModels) && serialModels[index].SerialNo == serial {
			continue
		}
		addUnitBySerial(a, serial)
	}
	log.Info().Msgf("[%s] finished scanning", a.Code())
}

// addUnitBySerial takes a serial number from some railway company and
// save it to the database if the serial number maps to a unit number.
func addUnitBySerial(a adapters.Adapter, serial string) {
	// add a unit serial record
	serialModel := models.SerialModel{
		Adapter:  a.Code(),
		SerialNo: serial,
	}
	info, err := a.Info(serial)
	if err == nil {
		serialModel.UnitNo, err = a.UnitNo(serial, info)
	}
	if err != nil || serialModel.UnitNo == "" {
		log.Debug().Msgf("[%s] %s -> %v", a.Code(), serial, err)
		return
	}
	serialModel.Operator, err = adapters.Operator(a, serial, info)
	if err != nil || serialModel.Operator == "" {
		serialModel.Operator = "?"
	}
	serialModel.Add()
	serialModel.AddTrainOperationLogs(info)
}
