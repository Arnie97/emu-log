package adapters_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
)

func ExampleBuildURL() {
	for _, testCase := range urlTestCases() {
		bureauCode, serial, url := testCase[0], testCase[1], testCase[2]
		b := adapters.MustGetBureauByCode(bureauCode)
		if urlBuilt := adapters.BuildURL(b, serial); urlBuilt != url {
			fmt.Println(urlBuilt)
			fmt.Println(url)
		}
	}
	// Output:
}

func ExampleParseURL() {
	for _, testCase := range urlTestCases() {
		bCode, serial, url := testCase[0], testCase[1], testCase[2]
		if b, s := adapters.ParseURL(url); s != serial || b.Code() != bCode {
			fmt.Println(b, s)
		}
	}

	fmt.Println(adapters.ParseURL("https://moerail.ml"))
	// Output: <nil>
}

func urlTestCases() (testCases [][]string) {
	content, err := ioutil.ReadFile("testdata/url.json")
	common.Must(err)
	common.Must(json.Unmarshal(content, &testCases))
	return
}

func mockHTTPClientRespBodyFromFile(mockFile string) {
	content, err := ioutil.ReadFile(filepath.Join("testdata", mockFile))
	common.Must(err)
	common.MockHTTPClientRespBody(string(content))
}

func printTrainNo(b adapters.Bureau, mockFiles ...string) {
	for _, mockFile := range mockFiles {
		mockHTTPClientRespBodyFromFile(mockFile)
		info, err := b.Info("")
		trainNo, date, err := b.TrainNo(info)
		fmt.Printf("%#-14v %-5v %#v\n", trainNo, err != nil, date)
	}
}

func printVehicleNo(b adapters.Bureau, mockFiles ...string) {
	for _, mockFile := range mockFiles {
		mockHTTPClientRespBodyFromFile(mockFile)
		info, err := b.Info("")
		vehicleNo, err := b.VehicleNo(info)
		fmt.Printf("%#-14v %-5v\n", vehicleNo, err != nil)
	}
}
