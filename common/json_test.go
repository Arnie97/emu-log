package common_test

import (
	"fmt"
	"os"

	"github.com/arnie97/emu-log/common"
)

func ExampleConf() {
	common.AppPath()
	file, err := os.Create(common.AppPath() + "/emu-log.json")
	if err != nil {
		fmt.Println(err)
	}
	file.Write([]byte(`{"a": "1234", "b": "5678"}`))
	file.Close()

	fmt.Println(common.Conf("a"), common.Conf("b"))
	// Output:
	// 1234 5678
}

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
	common.StructDecode(
		map[string]interface{}{"root": []float32{1, 2, 3}},
		&dest,
	)
	fmt.Printf("%+v", dest)
	// Output: {Field:[1 2 3]}
}
