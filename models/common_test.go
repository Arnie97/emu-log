package models_test

import (
	"fmt"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/arnie97/emu-log/models"
)

func ExampleDB() {
	x := models.DB()
	y := models.DB()
	_, err := x.Exec(`SELECT 1;`)
	c1 := models.CountRecords("emu_qrcode")
	c2 := models.CountRecords("emu_log", "date")

	fmt.Println("x == y:  ", x == y)
	fmt.Println("x != nil:", x != nil)
	fmt.Println("error:   ", err)
	fmt.Println("count:   ", c1, c2)
	// Output:
	// x == y:   true
	// x != nil: true
	// error:    <nil>
	// count:    0 0
}

func ExampleListLatestTrainByCondition() {
	resetTestDB()

	models.SerialModel{VehicleNo: "CRH6A4002", BureauCode: "N", SerialNo: "002"}.Add()
	fmt.Println(models.ListSerials(adapters.MustGetBureauByCode("N")))
	fmt.Println(models.ListSerialsForSingleVehicle("%J2015"))
	fmt.Println(models.ListLatestSerialForMultiVehicles(adapters.MustGetBureauByCode("P")))
	fmt.Println(models.ListVehiclesForSingleTrain("D5461"))
	models.LogModel{Date: "2020-11-20", VehicleNo: "CRH6A4002", TrainNo: "D5464/1/4"}.Add()
	fmt.Println(models.ListVehiclesForSingleTrain("D5461"))
	fmt.Println(models.ListLatestVehicleForMultiTrains([]string{"D5461", "D3045", "G666", "D5464"}))
	fmt.Println(models.ListTrainsForSingleVehicle("%2015"))
	fmt.Println(len(models.ListVehiclesForSingleTrain("D3071")))
	for i := 0; i < 10; i++ {
		models.LogModel{VehicleNo: "CRH2A2015", TrainNo: "D3074/D3071"}.Add()
	}
	fmt.Println(len(models.ListVehiclesForSingleTrain("D3071")))
	fmt.Println(models.ListLatestTrainForMultiVehicles([]string{"CRH6A4002", "CR200J2040", "CRH2A2460"}))
	fmt.Println(models.ListLatestTrainForMatchedVehicles("%A%02%"))

	// Output:
	// [{CRH6A4002 N 002} {CH001 N 053} {CRH2650 N 111} {CRH5A5075 N 472}]
	// [{CR200J2015 H PQ0916500} {CR200J2015 H PQ0916000}]
	// [{CR400AF0207 P 50704500} {CR400AF2015 P 50880000}]
	// [{2020-11-16 CR200J2015 D5464/1/4} {2020-11-14 CR200J2015 D5464/1/4} {2020-11-13 CR200J2040 D5464/1/4}]
	// [{2020-11-20 CRH6A4002 D5464/1/4} {2020-11-16 CR200J2015 D5464/1/4} {2020-11-14 CR200J2015 D5464/1/4} {2020-11-13 CR200J2040 D5464/1/4}]
	// [{2020-11-20 CRH6A4002 D5461} {2020-11-20 CRH6A4002 D5464} {2020-11-14 CR400AF0207 G666}]
	// [{2020-11-17 CR200J2015 D5468/D5465} {2020-11-16 CR200J2015 D5456/D5457} {2020-11-16 CR200J2015 D5464/1/4} {2020-11-15 CR200J2015 D5466/D5467} {2020-11-15 CR200J2015 D5468/D5465} {2020-11-15 CR200J2015 D5462/D5463} {2020-11-15 CR200J2015 D5458/D5455} {2020-11-14 CR200J2015 D5456/D5457} {2020-11-14 CR200J2015 D5464/1/4} {2020-11-13 CR200J2015 D5466/D5467} {2020-11-13 CR200J2015 D5468/D5465} {2020-10-26 CR400AF2015 G666} {2020-11-18 CRH2A2015 D3220} {2020-11-18 CRH2A2015 D3205} {2020-11-16 CRH2A2015 D3074/D3071} {2020-11-15 CRH2A2015 D3074/D3071} {2020-11-15 CRH2A2015 D3072/D3073} {2020-11-10 CRH2A2015 D3206} {2020-11-10 CRH2A2015 D3219}]
	// 2
	// 3
	// [{2020-11-13 CR200J2040 D5464/1/4} {2020-11-20 CRH6A4002 D5464/1/4}]
	// [{2020-11-15 CR400AF0207 G8907} {2020-11-20 CRH6A4002 D5464/1/4}]
}

func resetTestDB() {
	content := common.ReadMockFile("test.sql")
	_, err := models.DB().Exec(string(content))
	common.Must(err)
}
