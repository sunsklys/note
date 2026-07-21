package base

import (
	"reflect"
	"testing"
)

func TestMinimumRecolors(t *testing.T) {
	tests := []struct {
		str  string
		k    int
		want int
	}{
		{
			str:  "WBBWWBBWBW",
			k:    7,
			want: 3,
		},
	}

	for _, tt := range tests {
		got := minimumRecolors(tt.str, tt.k)
		if reflect.DeepEqual(got, tt.want) {
			t.Errorf("getAverages(%v, %v) = %v, want %v", tt.str, tt.k, got, tt.want)
		}
	}
}

// https://leetcode.cn/problems/minimum-recolors-to-get-k-consecutive-black-blocks/
func minimumRecolors(blocks string, k int) int {
	return 0
}
