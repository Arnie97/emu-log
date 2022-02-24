package tasks

import (
	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/models"
	"github.com/rs/zerolog/log"
)

// scanVehicleNo tries each unknown QR code in the brute force key space to see
// if any of these serial numbers was recently (or is currently) put in to use.
func scanVehicleNo(b adapters.Bureau) {
	log.Info().Msgf("[%s] started scanning for new vehicles", b.Code())
	defer wg.Done()

	serials := make(chan string)
	go func() {
		b.BruteForce(serials)
		close(serials)
	}()

	var index int
	serialModels := models.ListSerials(b)
	for serial := range serials {
		// skip existing serial numbers in the database
		for index < len(serialModels) && serialModels[index].SerialNo < serial {
			index++
		}
		if index < len(serialModels) && serialModels[index].SerialNo == serial {
			continue
		}
		addVehicleBySerial(b, serial)
	}
	log.Info().Msgf("[%s] finished scanning", b.Code())
}

// addVehicleBySerial takes a serial number from some railway company and
// save it to the database if the serial number maps to a vehicle number.
func addVehicleBySerial(b adapters.Bureau, serial string) {
	// add a vehicle serial record
	serialModel := models.SerialModel{
		BureauCode: b.Code(),
		SerialNo:   serial,
	}
	info, err := b.Info(serial)
	if err == nil {
		serialModel.VehicleNo, err = b.VehicleNo(serial, info)
	}
	if err != nil || serialModel.VehicleNo == "" {
		log.Debug().Msgf("[%s] %s -> %v", b.Code(), serial, err)
		return
	}
	serialModel.Add()
	serialModel.AddTrainOperationLogs(info)
}
