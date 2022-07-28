package common_test

import (
	"fmt"

	"github.com/arnie97/emu-log/common"
)

func ExamplePrettyPrint() {
	common.PrettyPrint(map[string]interface{}{
		"CIT380A": "CRH2C-2150",
		"CR200J":  nil,
		"CR400AF": []string{"0207", "0208"},
		"CR400BF": []string{"0503", "0507", "0305"},
		"CRH6A":   4002,
	})
	// Output:
	// {
	//     "CIT380A": "CRH2C-2150",
	//     "CR200J": null,
	//     "CR400AF": [
	//         "0207",
	//         "0208"
	//     ],
	//     "CR400BF": [
	//         "0503",
	//         "0507",
	//         "0305"
	//     ],
	//     "CRH6A": 4002
	// }
}

func ExampleGetField() {
	station := struct {
		Pinyin, Telegraphy string
		TMIS               int
	}{"HGT", "HTT", 53144}

	fmt.Println(
		common.GetField(station, "Pinyin"),
		common.GetField(station, "Telegraphy"),
		common.GetField(station, "TMIS"),
	)
	// Output: HGT HTT 53144

	defer func() {
		if recover() == nil {
			fmt.Println("panic expected here!")
		}
	}()
	common.GetField(station, "Nonexistent")
}

func ExampleStructDecode() {
	var dest struct {
		Field []int64 `json:"root"`
	}
	testCases := []interface{}{
		func() {},
		map[string]interface{}{"root": "123"},
		map[string]interface{}{"root": []float32{1, 2, 3}},
	}

	for _, testCase := range testCases {
		err := common.StructDecode(testCase, &dest)
		fmt.Printf("%+v %v\n", dest, err != nil)
	}
	// Output:
	// {Field:[]} true
	// {Field:[]} true
	// {Field:[1 2 3]} false
}
