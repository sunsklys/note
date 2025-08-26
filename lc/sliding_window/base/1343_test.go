package base

import "testing"

func TestNumOfSubarrays(t *testing.T) {
	tests := []struct {
		arr       []int
		k         int
		threshold int
		want      int
	}{
		{
			arr:       []int{2, 2, 2, 2, 5, 5, 5, 8},
			k:         3,
			threshold: 4,
			want:      3,
		},
		{
			arr:       []int{11, 13, 17, 23, 29, 31, 7, 5, 2, 3},
			k:         3,
			threshold: 5,
			want:      6,
		},
	}

	for _, tt := range tests {
		got := numOfSubarrays(tt.arr, tt.k, tt.threshold)
		if got != tt.want {
			t.Errorf("numOfSubarrays(%v, %d, %d) = %d; want %d", tt.arr, tt.k, tt.threshold, got, tt.want)
		}
	}
}

// https://leetcode.cn/problems/number-of-sub-arrays-of-size-k-and-average-greater-than-or-equal-to-threshold/
func numOfSubarrays(arr []int, k int, threshold int) int {
	result := 0
	w := 0
	for i, v := range arr {
		w += v
		if i < k-1 {
			continue
		}
		if w/k >= threshold {
			result++
		}
		out := arr[i-k+1]
		w -= out
	}
	return result
}
