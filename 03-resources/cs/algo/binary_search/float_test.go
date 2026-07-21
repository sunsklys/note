package binary_search

import (
	"testing"
)

func TestFloat(t *testing.T) {
	tests := []struct {
		name               string
		l, r, result, want float64
	}{
		{"用例", -100, 100, 1000.00, 9.99999999994543},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if r := float(v.l, v.r, v.result); r != v.want {
				t.Fatal("结果: ", r, "期望值: ", v.want)

			}
		})
	}
}

func float(l, r, x float64) float64 {
	for r-l > 1e-10 {
		mid := (r + l) / 2
		if x <= mid*mid*mid {
			r = mid
		} else {
			l = mid
		}
	}

	return l
}
