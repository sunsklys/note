package binary_search

import (
	"testing"
)

func TestRight(t *testing.T) {
	tests := []struct {
		name         string
		arr          []int
		search, want int
	}{
		{"一个", []int{1}, 1, 0},
		{"多个", []int{1, 2, 6, 9, 11}, 6, 2},
		{"两个重复", []int{1, 1}, 1, 1},
		{"多个重复", []int{1, 2, 6, 6, 9, 11}, 6, 3},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if r := right(v.arr, 0, len(v.arr)-1, v.search); r != v.want {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func right(s []int, l, r, x int) int {
	for l < r {
		mid := (l + r) + 1>>1
		if x >= s[mid] {
			l = mid
		} else {
			r = mid - 1
		}
	}

	return l
}
