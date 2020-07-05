package common_test

import (
	"testing"

	"github.com/arnie97/emu-log/common"
	"github.com/stretchr/testify/assert"
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

func TestGetField(t *testing.T) {
	station := struct {
		Pinyin, Telegraphy string
		TMIS               int
	}{"HGT", "HTT", 53144}

	assert.Equal(t, common.GetField(station, "Pinyin"), "HGT")
	assert.Equal(t, common.GetField(station, "Telegraphy"), "HTT")
	assert.Equal(t, common.GetField(station, "TMIS"), 53144)
	assert.Panics(t, func() {
		common.GetField(station, "Nonexistent")
	})
}
