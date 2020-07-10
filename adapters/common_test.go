package adapters_test

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
)

func mockHTTPClientRespBodyFromFile(mockFile string) {
	content, err := ioutil.ReadFile(filepath.Join("testdata", mockFile))
	common.Must(err)
	common.MockHTTPClientRespBody(string(content))
}

func printTrainNo(b adapters.Bureau, mockFile string) {
	mockHTTPClientRespBodyFromFile(mockFile)
	trainNo, date, err := b.TrainNo("")
	fmt.Printf("%#-14v %-5v %#v\n", trainNo, err != nil, date)
}

func printVehicleNo(b adapters.Bureau, mockFile string) {
	mockHTTPClientRespBodyFromFile(mockFile)
	vehicleNo, err := b.VehicleNo("")
	fmt.Printf("%#-14v %-5v\n", vehicleNo, err != nil)
}
