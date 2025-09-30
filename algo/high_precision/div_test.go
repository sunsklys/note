package high_precision

import (
	"strconv"
	"testing"
)

func TestDiv(t *testing.T) {
	tests := []struct {
		name, str1 string
		number     int
		quotient   string
		remainder  int
	}{
		{"整除", "100", 1, "100", 0},
		{"余数", "733", 10, "73", 3},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			if q, r := div(v.str1, v.number); q != v.quotient || r != v.remainder {
				t.Fatal("结果: ", q, r, "期望值: ", v.quotient, v.remainder)
			}
		})
	}
}

func div(s1 string, s2 int) (string, int) {
	t := 0
	s3 := make([]int, len(s1))
	for i := 0; i < len(s1); i++ {
		r := t*10 + int(s1[i]-'0')
		s3[i] = r / s2
		t = r % s2
	}

	for i := 0; i < len(s3); i++ {
		if s3[i] != 0 {
			s3 = s3[i:]
			break
		}
	}

	r := ""
	for i := 0; i < len(s3); i++ {
		r += strconv.Itoa(s3[i])
	}
	return r, t
}
