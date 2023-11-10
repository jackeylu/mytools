package util

import (
	"regexp"
	"strings"

	"github.com/yanyiwu/gojieba"
)

func IsAllCharacterDigit(str string) bool {
	for _, v := range str {
		if v < '0' || v > '9' {
			return false
		}
	}
	return true
}

// LongestCommonSubstr 最长公共子串
func LongestCommonSubstr(s1, s2 string) string {
	m := len(s1)
	n := len(s2)
	dp := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = make([]int, n+1)
	}

	maxLen := 0   // 最长公共子串的长度
	endIndex := 0 // 最长公共子串的结束索引

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > maxLen {
					maxLen = dp[i][j]
					endIndex = i
				}
			}
		}
	}

	return s1[endIndex-maxLen : endIndex]
}

func isChineseName(name string) bool {
	if name == "" || len(name) < 2 || len(name) > 12 {
		return false
	}

	x := gojieba.NewJieba()
	defer x.Free()

	// 判断是否只包含中文字符
	onlyChinese, _ := regexp.MatchString(`[\\u4e00-\\u9fa5]+`, name)
	if !onlyChinese {
		return false
	}

	// 判断是否包含其他非中文字符
	otherChars, _ := regexp.MatchString(`[^\\u4e00-\\u9fa5]+`, name)
	if otherChars {
		return false
	}

	// 判断是否包含英文字母
	englishLetters, _ := regexp.MatchString(`[a-zA-Z]+`, name)
	if englishLetters {
		return false
	}

	// 判断是否包含数字
	numbers, _ := regexp.MatchString(`[0-9]+`, name)
	if numbers {
		return false
	}

	// 判断是否包含空格
	if strings.Contains(name, " ") {
		return false
	}

	// 判断是否包含其他特殊字符
	specialChars := []string{"、", "（", "）"} // 根据具体要求添加特殊字符
	for _, char := range specialChars {
		if strings.Contains(name, char) {
			return false
		}
	}

	return true
}
