package sliding_window

import "testing"

func TestFindMaxAverage(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want float64
	}{
		{
			nums: []int{1, 12, -5, -6, 50, 3},
			k:    4,
			want: 12.75,
		},
		{
			nums: []int{5},
			k:    1,
			want: 5.00000,
		},
	}

	for _, tt := range tests {
		got := findMaxAverage(tt.nums, tt.k)
		if got != tt.want {
			t.Errorf("findMaxAverage(%v, %d) = %f, want %f", tt.nums, tt.k, got, tt.want)
		}
	}
}

// https://leetcode.cn/problems/maximum-average-subarray-i/description/
func findMaxAverage(nums []int, k int) float64 {
	w := 0
	result := float64(0)
	for i, v := range nums {
		w += v
		if i < k-1 {
			continue
		}

		result = max(result, float64(w)/float64(k))
		out := nums[i-k+1]
		w -= out
	}

	return result
}
