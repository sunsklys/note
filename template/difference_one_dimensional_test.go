package template

import (
	"reflect"
	"testing"
)

func TestDifferenceOneDimensional(t *testing.T) {
	tests := []struct {
		name      string
		arr       []int
		l, r, sum int
		want      []int
	}{
		{"第一位", []int{1, 2, 3, 4, 5}, 0, 1, 1, []int{2, 3, 3, 4, 5}},
		{"最后一位", []int{1, 2, 3, 4, 5}, 4, 4, 1, []int{1, 2, 3, 4, 6}},
		{"中间位", []int{1, 2, 3, 4, 5}, 1, 3, 1, []int{1, 3, 4, 5, 5}},
		{"所有", []int{1, 2, 3, 4, 5}, 0, 4, 1, []int{2, 3, 4, 5, 6}},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			r := differenceOneDimensional(v.arr, v.l, v.r, v.sum)
			if !reflect.DeepEqual(r, v.want) {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func differenceOneDimensional(s []int, l, r, sum int) []int {
	l = l + 1
	r = r + 1

	t := make([]int, len(s)+2)
	for i := 0; i < len(s); i++ {
		differenceOneDimensionalFormant(t, i+1, i+1, s[i])
	}

	differenceOneDimensionalFormant(t, l, r, sum)

	for i := 1; i < len(t); i++ {
		t[i] += t[i-1]
	}

	t = t[1 : len(t)-1]

	return t
}

func differenceOneDimensionalFormant(s []int, l, r, sum int) {
	s[l] += sum
	s[r+1] -= sum
}
