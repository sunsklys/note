package template

import (
	"reflect"
	"testing"
)

func TestSortMerge(t *testing.T) {
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
			sortMerge(v.arr, 0, len(v.arr)-1)
			if !reflect.DeepEqual(v.arr, v.want) {
				t.Fatal("结果: ", v.arr, "期望值: ", v.want)
			}
		})
	}
}

func sortMerge(s []int, l int, r int) {
	if l >= r {
		return
	}

	mid := (l + r) >> 1
	sortMerge(s, l, mid)
	sortMerge(s, mid+1, r)

	left, right := l, mid+1
	var t []int
	for left <= mid && right <= r {
		if s[left] <= s[right] {
			t = append(t, s[left])
			left++
		} else {
			t = append(t, s[right])
			right++
		}
	}

	for left <= mid {
		t = append(t, s[left])
		left++
	}

	for right <= r {
		t = append(t, s[right])
		right++
	}

	for left, right = l, 0; left <= r; left, right = left+1, right+1 {
		s[left] = t[right]
	}
}
