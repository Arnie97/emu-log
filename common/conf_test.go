package common_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/arnie97/emu-log/common"
	"github.com/pelletier/go-toml"
)

func TestRequestLimits(t *testing.T) {
	var conf common.GlobalConf
	err := toml.Unmarshal([]byte(`
[[limits]]
max-matches = 15
[limits.conditions]
User-Agent = "TrustedClient"
Authorization = "secret"

[[limits]]
max-matches = 300
[limits.conditions]
Authorization = "backup"
`), &conf)
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name        string
		headers     http.Header
		wantLimit   int
		wantMatched bool
	}{
		{"all headers match", http.Header{"User-Agent": {"TrustedClient/1.0"}, "Authorization": {"Bearer secret"}}, 15, true},
		{"second rule matches", http.Header{"Authorization": {"backup"}}, 300, true},
		{"first matching rule wins", http.Header{"User-Agent": {"TrustedClient/1.0"}, "Authorization": {"secret backup"}}, 15, true},
		{"one header missing", http.Header{"User-Agent": {"TrustedClient/1.0"}}, 0, false},
		{"one header does not match", http.Header{"User-Agent": {"OtherClient"}, "Authorization": {"secret"}}, 0, false},
		{"no headers", nil, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limit := conf.RequestLimit(tt.headers)
			if (limit != nil) != tt.wantMatched {
				t.Fatalf("RequestLimit() = %v, want matched: %v", limit, tt.wantMatched)
			}
			if limit != nil && limit.MaxMatches != tt.wantLimit {
				t.Errorf("RequestLimit().MaxMatches = %d, want %d", limit.MaxMatches, tt.wantLimit)
			}
		})
	}

	if limit := (&common.GlobalConf{}).RequestLimit(nil); limit != nil {
		t.Errorf("RequestLimit() with no rules = %v, want nil", limit)
	}
	if limit := (&common.GlobalConf{Limits: []common.LimitConf{{MaxMatches: 15}}}).RequestLimit(nil); limit == nil || limit.MaxMatches != 15 {
		t.Errorf("RequestLimit() with no conditions = %v, want a rule with MaxMatches 15", limit)
	}
}

func ExampleConf() {
	common.MockConf()

	fmt.Println(int64(common.Conf().Request.Interval))
	common.PrettyPrint(common.Conf().Request.Interval)

	serials := make(chan string)
	for _, a := range common.Conf().Adapters {
		if len(a.SearchSpace) > 0 {
			go a.EmitSerials(serials)
		}
	}
	for s := range serials {
		fmt.Println(s)
	}

	// Output:
	// 2460
	// "2.46µs"
	// CRH5-001A
	// CRH3-002C
	// CRH3-004C
	// CRH3-009C
	// CRH2-001A
	// CRH2-002A
	// CRH2-003A
}
