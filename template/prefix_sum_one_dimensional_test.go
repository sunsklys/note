package template

import (
	"testing"
)

func TestPrefixSumOneDimensional(t *testing.T) {
	tests := []struct {
		name       string
		arr        []int
		l, r, want int
	}{
		{"第一位", []int{1, 2, 3, 4, 5}, 0, 2, 6},
		{"所有", []int{1, 2, 3, 4, 5}, 0, 4, 15},
		{"中间", []int{1, 2, 3, 4, 5}, 2, 3, 7},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if r := prefixSumOneDimensional(v.arr, v.l, v.r); r != v.want {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func prefixSumOneDimensional(s []int, l, r int) int {
	l += 1
	r += 1

	s = append([]int{0}, s...)
	for i := 1; i < len(s); i++ {
		s[i] += s[i-1]
	}

	return s[r] - s[l-1]
}
