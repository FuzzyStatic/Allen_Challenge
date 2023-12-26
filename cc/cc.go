// Package helps validate credit card numbers
package cc

import (
	"regexp"
	"strings"
)

var (
	// Checks for 16 digits that start with 4, 5, or 6
	re16Digits = regexp.MustCompile(`^(?:4|5|6)[\d]{15,15}$`)

	// Checks for 16 digits that start with 4, 5, or 6 and has dashes between groups of 4 digits
	re16DigitsWithDashes = regexp.MustCompile(`^(?:4|5|6)[\d]{3,3}[-][\d]{4,4}[-][\d]{4,4}[-][\d]{4,4}$`)

	// Checks for a digit that repeats 4 or more times
	reRepeatingDigits = regexp.MustCompile(`0{4,}|1{4,}|2{4,}|3{4,}|4{4,}|5{4,}|6{4,}|7{4,}|8{4,}|9{4,}`)
)

func IsValid(ccNum string) bool {
	if !re16Digits.Match([]byte(ccNum)) && !re16DigitsWithDashes.Match([]byte(ccNum)) {
		return false
	}

	ccNumNoDashes := strings.Replace(ccNum, "-", "", -1)
	if reRepeatingDigits.Match([]byte(ccNumNoDashes)) {
		return false
	}

	return true
}
