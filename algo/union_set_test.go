package algo

import (
	"math"
	"testing"
)

func TestUnionSet(t *testing.T) {
	testes := []struct {
		name string
		set  [][]int
		want int
	}{
		{"一个集合", [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}}, 5},
		{"多个集合", [][]int{{1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}}, 7},
	}

	for _, test := range testes {
		t.Run(test.name, func(t *testing.T) {
			u := NewUnionSet(math.MaxInt8)
			for _, s := range test.set {
				u.Union(s[0], s[1])
			}
			if u.Find(1) != test.want {
				t.Errorf(" want: %d, got: %d", test.want, u.Find(1))
			}
		})
	}
}

type UnionSet struct {
	Father []int
}

func NewUnionSet(n int) *UnionSet {
	u := &UnionSet{
		Father: make([]int, n+1),
	}
	for i := 1; i <= n; i++ {
		u.Father[i] = i
	}

	return u
}

func (u *UnionSet) Union(i int, j int) {
	iFather := u.Father[i]
	jFather := u.Father[j]
	u.Father[iFather] = jFather
}

func (u *UnionSet) Find(i int) int {
	if i == u.Father[i] {
		return i
	} else {
		u.Father[i] = u.Find(u.Father[i])
		return u.Father[i]
	}
}
