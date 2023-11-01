package cmd

import (
	"testing"
)

func TestLotto(t *testing.T) {
	// 设置多个测试数据
	for _, tc := range []struct {
		input    int64
		expected int
	}{
		{5667, 4},
		{56678, 5},
		{566789, 6},
		{5667890, 7},
		{56678901, 8},
		{0, 1},
		{3, 1},
		{-1, 1},
		{1000000000000000000, 19},
	} {
		if i := getNumberLength(tc.input); i != tc.expected {
			t.Errorf("getNumberLength(%d) = %d, want %d", tc.input, i, tc.expected)
		}
	}

}
