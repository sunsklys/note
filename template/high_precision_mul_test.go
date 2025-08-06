package template

import (
	"slices"
	"strconv"
	"testing"
)

func TestHighPrecisionMul(t *testing.T) {
	tests := []struct {
		name, str1 string
		number     int
		want       string
	}{
		{"乘数单位被乘数单位", "1", 1, "1"},
		{"乘数多位被乘数单位", "111", 1, "111"},
		{"乘数单位被乘数多位", "6", 99, "594"},
		{"乘数多位被乘数多位", "1111", 99, "109989"},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if r := highPrecisionMul(v.str1, v.number); r != v.want {
				t.Fatal("结果: ", r, "期望值: ", v.want)
			}
		})
	}
}

func highPrecisionMul(str1 string, num int) string {
	rs1 := reverseMul(str1)

	var t int
	s3 := make([]int, len(rs1))
	for i := 0; i < len(rs1); i++ {
		t += rs1[i] * num
		s3[i] = t % 10
		t /= 10
	}

	for t > 0 {
		s3 = append(s3, t%10)
		t /= 10
	}

	slices.Reverse(s3)

	res := ""
	for i := 0; i < len(s3); i++ {
		res += strconv.Itoa(s3[i])
	}

	return res
}

func reverseMul(str string) []int {
	s := make([]int, len(str))
	for i := 0; i < len(str); i++ {
		s[i] = int(str[i] - '0')
	}
	slices.Reverse(s)
	return s
}
