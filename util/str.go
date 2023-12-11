package util

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
