package common

import (
	"regexp"
	"strings"
)

var (
	normalizer = strings.NewReplacer(
		"CRH400", "CR400",
		"CHR", "CRH",
		"1-", "",
		"-", "",
		"_", "",
	)
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

// ApproxEqualVehicleNo compares whether the proposed vehicle number
// is approximately the same as the original one.
func ApproxEqualVehicleNo(original, proposed string) bool {
	if strings.ContainsRune(proposed, '@') {
		return true
	}
	if len(original) > 4 {
		original = original[len(original)-4:]
	}
	if len(proposed) > 4 {
		proposed = proposed[len(proposed)-4:]
	}
	return original == proposed
}
