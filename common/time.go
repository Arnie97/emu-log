package common

import (
	"time"
)

const (
	// ISODate is the date layout defined in ISO 8601 or RFC 3339.
	ISODate = "2006-01-02"
)

var mockWallClock func() time.Time

// UnixMilli backports time.UnixMilli() from Go 1.17 and later versions.
func UnixMilli(args ...time.Time) int64 {
	var t time.Time

	switch {
	case mockWallClock != nil:
		t = mockWallClock()
	case len(args) == 0:
		t = time.Now()
	case len(args) == 1:
		t = args[0]
	default:
		panic("Invalid argument length")
	}

	return t.UnixNano() / 1e6
}

func MockStaticUnixMilli(t int64) {
	mockWallClock = func() time.Time {
		return time.Unix(t/1e3, t%1e3*1e6)
	}
}

// Duration is a TOML wrapper type for time.Duration.
// https://github.com/golang/go/issues/16039
type Duration time.Duration

// UnmarshalText parses a TOML string into a Duration value.
func (d *Duration) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		return nil
	}

	duration, err := time.ParseDuration(string(text))
	if err != nil {
		return err
	}

	*d = Duration(duration)
	return nil
}

// MarshalText formats a Duration value into a TOML string.
func (d Duration) MarshalText() (text []byte, err error) {
	return []byte(d.String()), nil
}

// String provide a human readable string representing the duration.
func (d Duration) String() string {
	return time.Duration(d).String()
}
