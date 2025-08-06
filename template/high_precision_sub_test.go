package template

import (
	"slices"
	"strconv"
	"testing"
)

func TestHighPrecisionSub(t *testing.T) {
	tests := []struct {
		name, str1, str2, want string
	}{
		{"单位不借位", "1", "1", "0"},
		{"单位不借位", "5", "10", "-5"},
		{"多位相同不借位", "111", "110", "1"},
		{"多位相同借位", "55", "37", "18"},
		{"多位第一位多不借位", "1236", "30", "1206"},
		{"多位第一位多借位", "1233", "877", "356"},
		{"多位第二位多不进位", "66", "166", "-100"},
		{"多位第二位多借位", "33", "3333", "-3300"},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if r := highPrecisionSub(v.str1, v.str2); r != v.want {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func highPrecisionSub(str1, str2 string) string {
	s1 := reverseSub(str1)
	s2 := reverseSub(str2)

	pre := ""
	if !compare(s1, s2) {
		s1, s2 = s2, s1
		pre = "-"
	}

	var t int
	s3 := make([]int, len(s1))
	for i := 0; i < len(s1); i++ {
		t = s1[i] - t
		if i < len(s2) {
			t -= s2[i]
		}

		s3[i] = (t + 10) % 10
		if t < 0 {
			t = 1
		} else {
			t = 0
		}
	}

	slices.Reverse(s3)

	for i := 0; i < len(s3); i++ {
		if s3[i] != 0 {
			s3 = s3[i:]
			break
		}
	}

	res := ""
	for i := 0; i < len(s3); i++ {
		res += strconv.Itoa(s3[i])
	}
	return pre + res
}

func reverseSub(str string) []int {
	s := make([]int, len(str))
	for i := 0; i < len(str); i++ {
		s[i] = int(str[i] - '0')
	}
	slices.Reverse(s)
	return s
}

func compare(a, b []int) bool {
	if len(a) != len(b) {
		return len(a) > len(b)
	}

	for i := len(a) - 1; i >= 0; i-- {
		if a[i] != b[i] {
			return a[i] > b[i]
		}
	}

	return true
}
