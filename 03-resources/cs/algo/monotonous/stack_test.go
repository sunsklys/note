package monotonous

import (
	"reflect"
	"testing"
)

func TestStack(t *testing.T) {
	tests := []struct {
		name string
		arr  []int
		want []int
	}{
		{"一个", []int{1}, []int{-1}},
		{"多个", []int{3, 4, 2, 7, 5}, []int{-1, 3, -1, 2, 2}},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if r := stack(v.arr); !reflect.DeepEqual(r, v.want) {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func stack(arr []int) []int {
	var res = make([]int, len(arr))
	var stack []int
	for i := 0; i < len(arr); i++ {
		for len(stack) > 0 && stack[len(stack)-1] >= arr[i] {
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			res[i] = -1
		} else {
			res[i] = stack[len(stack)-1]
		}

		stack = append(stack, arr[i])
	}
	return res
}
