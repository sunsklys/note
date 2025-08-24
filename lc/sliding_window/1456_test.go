package sliding_window

import "testing"

func TestNameMaxVowels(t *testing.T) {
	tests := []struct {
		s       string
		k, want int
	}{
		{"abciiidef", 3, 3},
		{"aeiou", 2, 2},
		{"leetcode", 3, 2},
		{"rhythms", 4, 0},
		{"tryhard", 4, 1},
	}

	for _, tt := range tests {
		if ans := maxVowels(tt.s, tt.k); ans != tt.want {
			t.Errorf("maxVowels(%s, %d) = %d; want %d", tt.s, tt.k, ans, tt.want)
		}
	}
}

// https://leetcode.cn/problems/maximum-number-of-vowels-in-a-substring-of-given-length/
func maxVowels(s string, k int) (ans int) {
	w := 0
	for i, in := range s {
		if in == 'a' || in == 'e' || in == 'i' || in == 'o' || in == 'u' {
			w++
		}
		if i < k-1 {
			continue
		}
		ans = max(ans, w)
		out := s[i-k+1]
		if out == 'a' || out == 'e' || out == 'i' || out == 'o' || out == 'u' {
			w--
		}
	}
	return
}
