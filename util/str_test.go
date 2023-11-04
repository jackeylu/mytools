package util

import (
	"fmt"
	"regexp"
	"testing"
)

func TestIsAllCharacterDigit(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool
	}{
		{"12345", true},
		{"abc123", false},
		{"987654321", true},
		{"", true},
	}

	for _, testCase := range testCases {
		result := IsAllCharacterDigit(testCase.input)
		if result != testCase.expected {
			t.Errorf("Test case failed: input = %s, expected = %t, got = %t", testCase.input, testCase.expected, result)
		}
	}
}

func TestReg(t *testing.T) {
	f := func(str string) {
		// str := "李四2200023011张三"
		regex := "(\u4e00-\u9fa5)(\\d+)"
		match, _ := regexp.Compile(regex)
		result := match.FindAllString(str, -1)
		fmt.Println(result)
	}

	f("李四2200023011张三")

}
