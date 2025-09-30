package base

import (
	"reflect"
	"testing"
)

func TestGetAverAges(t *testing.T) {
	tests := []struct {
		nums []int
		k    int
		want []int
	}{
		{
			nums: []int{7, 4, 3, 9, 1, 8, 5, 2, 6},
			k:    3,
			want: []int{-1, -1, -1, 5, 4, 4, -1, -1, -1},
		},
	}

	for _, tt := range tests {
		got := getAverages(tt.nums, tt.k)
		if reflect.DeepEqual(got, tt.want) {
			t.Errorf("getAverages(%v, %v) = %v, want %v", tt.nums, tt.k, got, tt.want)
		}
	}
}

// https://leetcode.cn/problems/k-radius-subarray-averages/
func getAverages(nums []int, k int) []int {
	result := make([]int, 0)
	sum := 0
	for i, v := range nums {
		sum += v
		if i-k < 0 {
			result = append(result, -1)
			continue
		}

		if i+k >= len(nums) {
			result = append(result, -1)
			continue
		}

		tSum := sum
		for j := i + 1; j < i+k+1; j++ {
			tSum += nums[j]
		}

		result = append(result, tSum/((2*k)+1))

		if i-k >= 0 {
			sum -= nums[i-k]
		}
	}

	return result
}
