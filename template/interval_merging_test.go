package template

import (
	"reflect"
	"sort"
	"testing"
)

func TestIntervalMerging(t *testing.T) {
	tests := []struct {
		name string
		arr  [][]int
		want [][]int
	}{
		{"一个", [][]int{{1, 2}}, [][]int{{1, 2}}},
		{"多个", [][]int{{1, 2}, {2, 4}, {5, 6}, {7, 8}, {7, 9}}, [][]int{{1, 4}, {5, 6}, {7, 9}}},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if r := intervalMerging(v.arr); !reflect.DeepEqual(r, v.want) {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func intervalMerging(s [][]int) [][]int {
	sort.Slice(s, func(i, j int) bool {
		return s[i][0] < s[j][0]
	})

	var res = make([][]int, 0)
	l := s[0][0]
	r := s[0][1]
	for i := 1; i < len(s); i++ {
		if s[i][0] > r {
			res = append(res, []int{l, r})
			l = s[i][0]
			r = s[i][1]
		} else {
			if r < s[i][1] {
				r = s[i][1]
			}
		}
	}

	res = append(res, []int{l, r})
	return res
}
