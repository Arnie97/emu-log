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
		if b, s := adapters.ParseURL(url); s != serial || b.Code() != bCode {
			fmt.Println(b, s)
		}
	}

	fmt.Println(adapters.ParseURL("https://moerail.ml"))
	// Output: <nil>
}

func urlTestCases() (testCases [][]string) {
	common.Must(json.Unmarshal(common.ReadMockFile("url.json"), &testCases))
	return
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
		info, err := b.Info("")
		trainNo, date, err := b.TrainNo(info)
		fmt.Printf("%#-14v %-5v %#v\n", trainNo, err != nil, date)
	}

	for _, mockBody := range []string{"", "null", "<html>not json</html>"} {
		common.MockHTTPClientRespBody(mockBody)
		info, err := b.Info("")
		if info != nil && err == nil {
			fmt.Printf("uncaught error for http response %#v", mockBody)
		}
	}
}

func printVehicleNo(b adapters.Bureau, mockFiles ...string) {
	for _, mockFile := range mockFiles {
		common.MockHTTPClientRespBodyFromFile(mockFile)
		info, err := b.Info("")
		vehicleNo, err := b.VehicleNo(info)
		fmt.Printf("%#-14v %-5v\n", vehicleNo, err != nil)
	}
}
