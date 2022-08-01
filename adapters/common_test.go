package adapters_test

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/arnie97/emu-log/adapters"
	"github.com/arnie97/emu-log/common"
	"github.com/pelletier/go-toml"
)

type (
	AdapterTestFile struct {
		Adapters map[string]AdapterTestDefinition `toml:"adapters"`
	}
	AdapterTestDefinition struct {
		SerialNoPattern string            `toml:"pattern"`
		TestCases       []AdapterTestCase `toml:"cases"`
	}
	AdapterTestCase struct {
		SerialNo string `toml:"serial"`
		URL      string `toml:"url"`
	}
)

func ExampleSessionID() {
	req, _ := http.NewRequest(http.MethodGet, "", nil)
	for _, b := range adapters.Bureaus {
		if transport, ok := b.(http.RoundTripper); ok {
			transport.RoundTrip(req)
		}
	}
	// Output:
}

func ExampleBuildURL() {
	for bureauCode, testDef := range getTests() {
		b := adapters.MustGetBureauByCode(bureauCode)
		item := testDef.TestCases[0]
		if urlBuilt := adapters.BuildURL(b, item.SerialNo); urlBuilt != item.URL {
			fmt.Println(urlBuilt)
			fmt.Println(item.URL)
		}
	}
	// Output:
}

func ExampleParseURL() {
	for bureauCode, testDef := range getTests() {
		for _, item := range testDef.TestCases {
			b, s := adapters.ParseURL(item.URL)
			if b == nil {
				fmt.Println(item.URL, "->", "?")
				continue
			}
			if b.Code() != bureauCode || s != item.SerialNo {
				fmt.Println(item.URL, "->", b.Name(), s)
			}
		}
	}

	fmt.Println(adapters.ParseURL("https://moerail.ml"))
	// Output: <nil>
}

func getTests() map[string]AdapterTestDefinition {
	var testFile AdapterTestFile
	common.Must(toml.Unmarshal(common.ReadMockFile("adapters.toml"), &testFile))
	return testFile.Adapters
}

func getMockSerialNo(b adapters.Bureau) string {
	return getTests()[b.Code()].TestCases[0].SerialNo
}

func assertBruteForce(b adapters.Bureau, assert func(string) bool) {
	b.AlwaysOn()
	serials := make(chan string, 1024)
	go func() {
		for _, rule := range adapters.AdapterConf(b).SearchSpace {
			rule.Emit(serials)
		}
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
		serialNo := getMockSerialNo(b)
		info, err := b.Info(serialNo)
		vehicleNo, err := b.VehicleNo(serialNo, info)
		fmt.Printf("%#-14v %v\n", vehicleNo, err != nil)
	}
}
