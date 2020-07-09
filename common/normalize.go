package common

import (
	"regexp"
	"strings"
)

const (
	ISODate = "2006-01-02"
)

var (
	normalizer    = strings.NewReplacer("-", "", "_", "")
	trainNoRegExp = regexp.MustCompile(`\b[GDC]?\d{1,4}\b`)
)

// NormalizeTrainNo converts possibly abbreviated train number pairs to an
// array of full qualified train number strings.
func NormalizeTrainNo(trainNo string) (results []string) {
	var initial string
	for i, part := range strings.Split(trainNo, "/") {
		if part = trainNoRegExp.FindString(part); len(part) == 0 {
			return
		} else if i == 0 && part[0] <= '9' {
			return
		} else if i == 0 {
			initial = part
		} else if omitted := len(initial) - len(part); omitted > 0 {
			part = initial[:omitted] + part
		}
		results = append(results, part)
	}
	return
}

func NormalizeVehicleNo(vehicleNo string) string {
	return normalizer.Replace(vehicleNo)
}
