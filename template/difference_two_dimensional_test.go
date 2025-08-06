package template

import (
	"reflect"
	"testing"
)

func TestDifferenceTwoDimensional(t *testing.T) {
	tests := []struct {
		name                string
		arr                 [][]int
		x1, y1, x2, y2, sum int
		want                [][]int
	}{
		{"多维",
			[][]int{
				{1, 2, 2, 1},
				{3, 2, 2, 1},
				{1, 1, 1, 1},
			},
			1, 1, 2, 2, 1,
			[][]int{
				{1, 2, 2, 1},
				{3, 3, 3, 1},
				{1, 2, 2, 1},
			},
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			r := differenceTwoDimensional(v.arr, v.x1, v.y1, v.x2, v.y2, v.sum)
			if !reflect.DeepEqual(r, v.want) {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func differenceTwoDimensional(s [][]int, x1, y1, x2, y2, sum int) [][]int {
	x1 += 1
	y1 += 1
	x2 += 1
	y2 += 1

	t := make([][]int, len(s)+2)
	for i := 0; i < len(t); i++ {
		t[i] = make([]int, len(s[0])+2)
	}

	for i := 0; i < len(s); i++ {
		for j := 0; j < len(s[i]); j++ {
			differenceTwoDimensionalFormat(t, i+1, j+1, i+1, j+1, s[i][j])
		}
	}

	differenceTwoDimensionalFormat(t, x1, y1, x2, y2, sum)

	for i := 1; i < len(t); i++ {
		for j := 1; j < len(t[i]); j++ {
			t[i][j] += t[i-1][j] + t[i][j-1] - t[i-1][j-1]
		}
	}

	for i := 0; i < len(t); i++ {
		t[i] = t[i][1 : len(t[i])-1]
	}

	t = t[1 : len(t)-1]

	return t
}

func differenceTwoDimensionalFormat(s [][]int, x1, y1, x2, y2, c int) {
	s[x1][y1] += c
	s[x2+1][y1] -= c
	s[x1][y2+1] -= c
	s[x2+1][y2+1] += c
}
