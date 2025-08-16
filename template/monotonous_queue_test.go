package template

import (
	"reflect"
	"testing"
)

func TestMonotonousQueue(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		siz  int
		want []int
	}{
		{"多个", []int{1, 3, -1, -3, 5, 3, 6, 7}, 3, []int{-1, -3, -3, -3, 3, 3}},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if r := monotonousQueue(v.arr, v.siz); !reflect.DeepEqual(r, v.want) {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func monotonousQueue(arr []int, size int) []int {
	var res = make([]int, 0)
	var queue []int
	for i := 0; i < len(arr); i++ {
		for len(queue) > 0 && queue[0] < i-size+1 {
			queue = queue[1:]
		}

		for len(queue) > 0 && arr[queue[len(queue)-1]] > arr[i] {
			queue = queue[:len(queue)-1]
		}

		queue = append(queue, i)
		if i >= size-1 {
			res = append(res, arr[queue[0]])
		}

	}
	return res
}
