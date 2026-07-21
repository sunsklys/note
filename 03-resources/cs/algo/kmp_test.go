package algo

import (
	"reflect"
	"testing"
)

func TestKmp(t *testing.T) {
	tests := []struct {
		name       string
		str1, str2 string
		want       []int
	}{
		{"正常", "aba", "ababa", []int{0, 2}},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if r := KMP(v.str1, v.str2); !reflect.DeepEqual(r, v.want) {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func KMP(str1, str2 string) []int {
	str1 = " " + str1
	str2 = " " + str2
	res := make([]int, 0)
	var next = make([]int, len(str1))
	for i, j := 2, 0; i < len(str1); i++ {
		for j > 0 && str1[i] != str1[j+1] {
			j = next[j]
		}

		if str1[i] == str1[j+1] {
			j++
		}

		next[i] = j
	}

	for i, j := 1, 0; i < len(str2); i++ {
		for j > 0 && str2[i] != str1[j+1] {
			j = next[j]
		}

		if str2[i] == str1[j+1] {
			j++
		}

		if j == len(str1)-1 {
			res = append(res, i-(len(str1)-1))
			j = next[j]
		}
	}

	return res
}
