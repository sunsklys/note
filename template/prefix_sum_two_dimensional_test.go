package template

import (
	"testing"
)

func TestPrefixSumTwoDimensional(t *testing.T) {
	tests := []struct {
		name                 string
		arr                  [][]int
		x1, y1, x2, y2, want int
	}{
		{"一维", [][]int{{1, 7, 2, 4}}, 0, 0, 0, 1, 8},
		{"多维", [][]int{{1, 7, 2, 4}, {3, 6, 2, 8}, {2, 1, 2, 3}}, 0, 0, 1, 1, 17},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if r := prefixSumTwoDimensional(v.arr, v.x1, v.y1, v.x2, v.y2); r != v.want {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func prefixSumTwoDimensional(s [][]int, x1, y1, x2, y2 int) int {
	x1 += 1
	y1 += 1
	x2 += 1
	y2 += 1

	t := make([][]int, len(s))
	for i := 0; i < len(s); i++ {
		t[i] = append([]int{0}, s[i]...)
	}

	sum := append([][]int{make([]int, len(t[0]))}, t...)
	for i := 1; i < len(sum); i++ {
		for j := 1; j < len(sum[i]); j++ {
			sum[i][j] += sum[i-1][j] + sum[i][j-1] - sum[i-1][j-1]
		}
	}

	return sum[x2][y2] - sum[x1-1][y2] - sum[x2][y1-1] + sum[x1-1][y1-1]
}
