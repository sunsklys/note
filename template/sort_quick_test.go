package template

import (
	"reflect"
	"testing"
)

func TestSortQuick(t *testing.T) {
	tests := []struct {
		name      string
		arr, want []int
	}{
		{"一个", []int{1}, []int{1}},
		{"两个不同", []int{1, 3}, []int{1, 3}},
		{"两个相同", []int{1, 1}, []int{1, 1}},
		{"多个不同", []int{1, 4, 3, 6, 5}, []int{1, 3, 4, 5, 6}},
		{"多个相同", []int{1, 4, 3, 4, 5}, []int{1, 3, 4, 4, 5}},
		{"一个空切片", []int{}, []int{}},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			sortQuick(v.arr, 0, len(v.arr)-1)
			if !reflect.DeepEqual(v.arr, v.want) {
				t.Fatal("结果: ", v.arr, "期望值: ", v.want)
			}
		})
	}
}

func sortQuick(s []int, l int, r int) {
	if l >= r {
		return
	}

	left, right, mid := l-1, r+1, s[(l+r)>>1]
	for left < right {
		for {
			left++
			if s[left] >= mid {
				break
			}
		}

		for {
			right--
			if s[right] <= mid {
				break
			}
		}

		if left < right {
			s[left], s[right] = s[right], s[left]
		}
	}

	sortQuick(s, l, right)
	sortQuick(s, right+1, r)
}
