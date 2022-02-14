package common

import (
	"regexp"
	"strings"
)

var (
	vehicleNoNormalizeRules = []struct{ pattern, replace string }{
		{`[1-4]-`, ""},
		{"-", ""},
		{"_", ""},
		{"(CRH380D)V", "$1"},
		{"(CR)H([34]00)", "$1$2"},
		{"CHR", "CRH"},
	}
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
	for _, rule := range vehicleNoNormalizeRules {
		vehicleNo = regexp.MustCompile(
			rule.pattern).ReplaceAllString(vehicleNo, rule.replace)
	}
	return vehicleNo
}

// ApproxEqualVehicleNo compares whether the proposed vehicle number
// is approximately the same as the original one.
func ApproxEqualVehicleNo(original, proposed string) bool {
	if len(original) == 0 || strings.ContainsRune(proposed, '@') {
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
