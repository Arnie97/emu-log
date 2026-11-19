package adapters_test

import (
	"testing"

	"github.com/arnie97/emu-log/adapters"
)

func TestOfflinePSRAdapterParseBuildURL(t *testing.T) {
	for _, input := range []string{
		"https://p.12306.cn/tservice/qr/offlinePSR/v1?100S0**ztijYx9NoSVfSy8gOR+zRHpYoY5lE1YBQzTvIiTiIcckw4ARkDbss4WGjuvB/NM/XCrBOqMBwZYaL/Yc86VtFD/6J/D0lzGGL0emox9JRiQBXRCDtgFy1sIVwXC8yfFhPSggUAyJQr/9Z2jusioldLpASXvnPaGxHh/+qWyFIdas4TisXcU5MYPyJKwgpaEdqX2fXQ==",
		"https://p.12306.cn/tservice/qr/offlinePSR/v1?110S0**tzqZs9PbCctl6Y1dHFxtF7RnIcUzSX/oKbYvEN976bq3f3kDv5FUA84G5dQppwTw0n3kSGXYIpQFb3GUIMaHZwkmg4Hu9ErAWTMJId9rGdIYu6mZQkL1Yz7ga7hC5pXfCZpSiQTIaF2Toh63SkLVF/Oe5BBcBa/cgWm+WM91ywb3TOV7w3BcSmX6tiIxZiwrPkgnKg==",
	} {
		var o adapters.Adapter = adapters.OfflinePSR{}
		a, serial := adapters.ParseURL(input)
		if a == nil || a.Code() != o.Code() {
			t.Fatal("ParseURL: provided URL was not recognized by the offlinePSR adapter")
		}
		info, err := a.Info(serial)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := a.UnitNo(serial, info); err != nil {
			t.Fatal(err)
		}
		if adapters.BuildURL(a, serial) != input {
			t.Error("BuildURL did not reproduce the provided offlinePSR URL")
		}
	}
}
