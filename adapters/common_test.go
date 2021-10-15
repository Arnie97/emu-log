package adapters_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
)

func ExampleRoundTripper() {
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	for _, b := range adapters.Bureaus {
		if transport, ok := b.(http.RoundTripper); ok {
			transport.RoundTrip(req)
		}
	}
	// Output:
}

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
		b, s := adapters.ParseURL(url)
		if b == nil {
			fmt.Println(url, "->", "?")
			continue
		}
		if b.Code() != bCode || s != serial {
			fmt.Println(url, "->", b.Name(), s)
		}
	}

	fmt.Println(adapters.ParseURL("https://moerail.ml"))
	// Output: <nil>
}

func urlTestCases() (testCases [][]string) {
	common.Must(json.Unmarshal(common.ReadMockFile("url.json"), &testCases))
	return
}

func getMockSerialNo(b adapters.Bureau) string {
	for _, testCase := range urlTestCases() {
		bureauCode, serial, _ := testCase[0], testCase[1], testCase[2]
		if bureauCode == b.Code() {
			return serial
		}
	}
	return ""
}

func assertBruteForce(b adapters.Bureau, assert func(string) bool) {
	b.AlwaysOn()
	serials := make(chan string, 1024)
	go func() {
		b.BruteForce(serials)
		close(serials)
	}()
	for s := range serials {
		if !assert(s) {
			fmt.Printf("[%s] invalid serial number pattern: %s\n", b.Code(), s)
		}
	}
}

func assertBruteForceRegExp(b adapters.Bureau, pattern string) {
	assertBruteForce(b, regexp.MustCompile(pattern).MatchString)
}

func printTrainNo(b adapters.Bureau, mockFiles ...string) {
	b.Name()

	for _, mockFile := range mockFiles {
		common.MockHTTPClientRespBodyFromFile(mockFile)
		info, err := b.Info(getMockSerialNo(b))
		trains, err := b.TrainNo(info)
		fmt.Printf("\n%v\n", err != nil)
		for _, train := range trains {
			fmt.Printf("%#-14v %#v\n", train.TrainNo, train.Date)
		}
	}

	for _, mockBody := range []string{"", "null", "<html>not json</html>"} {
		common.MockHTTPClientRespBody(mockBody)
		info, err := b.Info(getMockSerialNo(b))
		if info != nil && err == nil {
			fmt.Printf("uncaught error for http response %#v", mockBody)
		}
	}

	common.MockHTTPClientError(fmt.Errorf("mock http error"))
	info, err := b.Info(getMockSerialNo(b))
	if info != nil && err == nil {
		fmt.Printf("uncaught error for http error")
	}
}

func printVehicleNo(b adapters.Bureau, mockFiles ...string) {
	for _, mockFile := range mockFiles {
		common.MockHTTPClientRespBodyFromFile(mockFile)
		info, err := b.Info(getMockSerialNo(b))
		vehicleNo, err := b.VehicleNo(info)
		fmt.Printf("%#-14v %v\n", vehicleNo, err != nil)
	}
}
