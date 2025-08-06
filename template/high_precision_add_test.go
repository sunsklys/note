package template

import (
	"slices"
	"strconv"
	"testing"
)

func TestHighPrecisionAdd(t *testing.T) {
	tests := []struct {
		name, str1, str2, want string
	}{
		{"单位不进位", "1", "1", "2"},
		{"单位进位", "5", "5", "10"},
		{"多位相同不进位", "11", "12", "23"},
		{"多位相同进位", "55", "55", "110"},
		{"多位第一位多不进位", "1236", "30", "1266"},
		{"多位第一位多进位", "1233", "877", "2110"},
		{"多位第二位多不进位", "66", "100", "166"},
		{"多位第二位多进位", "777", "5333", "6110"},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if r := highPrecisionAdd(v.str1, v.str2); r != v.want {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func highPrecisionAdd(str1 string, str2 string) string {
	if len(str1) < len(str2) {
		str1, str2 = str2, str1
	}

	s1 := reverseAdd(str1)
	s2 := reverseAdd(str2)

	var t int
	var s3 []int
	for i := 0; i < len(s1); i++ {
		t += s1[i]
		if i < len(s2) {
			t += s2[i]
		}

		s3 = append(s3, t%10)
		t /= 10
	}

	if t > 0 {
		s3 = append(s3, t%10)
	}

	slices.Reverse(s3)

	res := ""
	for i := 0; i < len(s3); i++ {
		res += strconv.Itoa(s3[i])
	}

	return res
}

func reverseAdd(str string) []int {
	s := make([]int, len(str))
	for i := 0; i < len(str); i++ {
		s[i] = int(str[i] - '0')
	}
	slices.Reverse(s)
	return s
}
