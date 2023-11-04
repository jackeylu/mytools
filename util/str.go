package util

func IsAllCharacterDigit(str string) bool {
	for _, v := range str {
		if v < '0' || v > '9' {
			return false
		}
	}
	return true
}
