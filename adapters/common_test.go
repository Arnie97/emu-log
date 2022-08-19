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
	MockAdapter struct {
		adapters.Shanghai
		code string
	}
)

func (m *MockAdapter) Code() string {
	return m.code
}

func ExampleAdapterConf() {
	common.MockConf()
	conf := adapters.AdapterConf(&MockAdapter{code: "X"})
	fmt.Println(conf.Request.UserAgent, int64(conf.Request.Interval))

	// Output:
	// Mozilla/5.0 4002
}

func ExampleSessionID() {
	common.MockConf()
	fmt.Println(adapters.SessionID(&MockAdapter{code: "X"}))

	req, _ := http.NewRequest(http.MethodGet, "", nil)
	for _, a := range adapters.Adapters {
		if transport, ok := a.(http.RoundTripper); ok {
			transport.RoundTrip(req)
		}
	}
	// Output:
	// hello-world
}

func ExampleBuildURL() {
	for adapterCode, testDef := range getTests() {
		a := adapters.MustGetAdapterByCode(adapterCode)
		item := testDef.TestCases[0]
		if urlBuilt := adapters.BuildURL(a, item.SerialNo); urlBuilt != item.URL {
			fmt.Println(urlBuilt)
			fmt.Println(item.URL)
		}
	}
	// Output:
}

func ExampleParseURL() {
	for adapterCode, testDef := range getTests() {
		for _, item := range testDef.TestCases {
			a, s := adapters.ParseURL(item.URL)
			if a == nil {
				fmt.Println(item.URL, "->", "?")
				continue
			}
			if a.Code() != adapterCode || s != item.SerialNo {
				fmt.Println(item.URL, "->", a.Name(), s)
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

func getMockSerialNo(a adapters.Adapter, preferredIndex int) string {
	testCases := getTests()[a.Code()].TestCases
	if len(testCases) > preferredIndex {
		return testCases[preferredIndex].SerialNo
	}
	return testCases[0].SerialNo
}

func assertBruteForce(a adapters.Adapter, assert func(string) bool) {
	a.AlwaysOn()
	serials := make(chan string, 1024)
	go func() {
		for _, rule := range adapters.AdapterConf(a).SearchSpace {
			rule.Emit(serials)
		}
		close(serials)
	}()
	for s := range serials {
		if !assert(s) {
			fmt.Printf("[%s] invalid serial number pattern: %s\n", a.Code(), s)
		}
	}
}

func assertBruteForceRegExp(a adapters.Adapter, pattern string) {
	assertBruteForce(a, regexp.MustCompile(pattern).MatchString)
}

func printTrainNo(a adapters.Adapter, mockFiles ...string) {
	a.Name()

	for index, mockFile := range mockFiles {
		common.MockHTTPClientRespBodyFromFile(mockFile)
		info, err := a.Info(getMockSerialNo(a, index))
		trains, err := a.TrainNo(info)
		fmt.Printf("\n%v\n", err != nil)
		for _, train := range trains {
			fmt.Printf("%#-14v %#v\n", train.TrainNo, train.Date)
		}
	}

	for _, mockBody := range []string{"", "null", "<html>not json</html>"} {
		common.MockHTTPClientRespBody(mockBody)
		info, err := a.Info(getMockSerialNo(a, 0))
		if info != nil && err == nil {
			fmt.Printf("uncaught error for http response %#v", mockBody)
		}
	}

	common.MockHTTPClientError(fmt.Errorf("mock http error"))
	info, err := a.Info(getMockSerialNo(a, 0))
	if info != nil && err == nil {
		fmt.Printf("uncaught error for http error")
	}
}

func printUnitNo(a adapters.Adapter, mockFiles ...string) {
	for index, mockFile := range mockFiles {
		common.MockHTTPClientRespBodyFromFile(mockFile)
		serialNo := getMockSerialNo(a, index)
		info, infoErr := a.Info(serialNo)
		unitNo, unitErr := a.UnitNo(serialNo, info)
		operator, err := adapters.Operator(a, serialNo, info)
		fmt.Printf("%1v %#-14v %5v %5v %5v\n", operator, unitNo,
			infoErr != nil, unitErr != nil, err != nil)
	}
}
