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

func TestLongestCommonSubstr(t *testing.T) {
	testCases := []struct {
		desc string
		s1   string
		s2   string
		want string
	}{
		{
			desc: "英文",
			s1:   "OldSite: The old URL of this website",
			s2:   "NewSite: The new URL of this website",
			want: " URL of this website",
		},
		{
			desc: "中文",
			s1:   "Lab5-文件上传+文件投票+目录遍历",
			s2:   "220301049李欣蕊-文件上传+文件投票+目录遍历",
			want: "-文件上传+文件投票+目录遍历",
		},
		{
			desc: "没有公共子串",
			s1:   "OldSite: The old URL of this website",
			s2:   "220301049李欣蕊-文件上传+文件投票+目录遍历",
			want: "",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := LongestCommonSubstr(tC.s1, tC.s2); got != tC.want {
				t.Errorf("got %q, want %q", got, tC.want)
			}
		})
	}
}
