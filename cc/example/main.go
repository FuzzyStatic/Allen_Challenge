package main

import (
	"Allen_Challenge/cc"
	"fmt"
)

// Credit card numbers to check, add whatever you like
var ccNums = []string{
	"6",
	"4123456789123456",
	"5123-4567-8912-3456",
	"61234-567-8912-3456",
	"5100-0067-8912-3456",
	"5111-1167-8912-3456",
	"5123456789123456",
	"5122-2267-8912-3456",
	"5133-3367-8912-3456",
	"5144-4467-8912-3456",
	"5155-5567-8912-3456",
	"5166-6667-8912-3456",
	"6123456789123456",
	"5177-7767-8912-3456",
	"5188-8867-8912-3456",
	"5199-9967-8912-3456",
	"5123 - 3567 - 8912 - 3456",
	"5133-336789123456",
	"6843728",
	"foobar",
}

func main() {
	for _, ccNum := range ccNums {
		if cc.IsValid(ccNum) {
			fmt.Println("Valid")
			continue
		}

		fmt.Println("Invalid")
	}
}
