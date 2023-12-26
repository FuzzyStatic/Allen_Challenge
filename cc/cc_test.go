package cc

import (
	"fmt"
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		ccNum            string
		expectedValidity bool
	}{
		{"6", false},
		{"4123456789123456", true},
		{"5123456789123456", true},
		{"6123456789123456", true},
		{"5123-4567-8912-3456", true},
		{"61234-567-8912-3456", false},
		{"5100-0067-8912-3456", false},
		{"5111-1167-8912-3456", false},
		{"5122-2267-8912-3456", false},
		{"5133-3367-8912-3456", false},
		{"5144-4467-8912-3456", false},
		{"5155-5567-8912-3456", false},
		{"5166-6667-8912-3456", false},
		{"5177-7767-8912-3456", false},
		{"5188-8867-8912-3456", false},
		{"5199-9967-8912-3456", false},
		{"5123 - 3567 - 8912 - 3456", false},
		{"5133-336789123456", false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v validity is %v", test.ccNum, test.expectedValidity), func(t *testing.T) {
			if IsValid(test.ccNum) != test.expectedValidity {
				t.Errorf("%v validity should not be %v", test.ccNum, test.expectedValidity)
			}
		})
	}
}
